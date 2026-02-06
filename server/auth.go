package server

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KaraberusClaims struct {
	IsAdmin bool `json:"is_admin"`
	Scopes
	oidc.IDTokenClaims
	jwt.MapClaims `json:"-"`
}

var (
	sessionCookieName = "karaberus_session"
	currentUserCtxKey = "current_user"
)

func authState() string {
	return uuid.New().String()
}

func addOidcRoutes(ctx context.Context, app *fiber.App) {
	provider, err := rp.NewRelyingPartyOIDC(
		ctx,
		CONFIG.OIDC.Issuer,
		CONFIG.OIDC.ClientID,
		CONFIG.OIDC.ClientSecret,
		fmt.Sprintf("%v/api/oidc/callback", CONFIG.Listen.BaseURL),
		CONFIG.OIDC.Scopes,
	)

	if err != nil {
		panic(err)
	}

	app.Get("/api/oidc/login", adaptor.HTTPHandler(rp.AuthURLHandler(authState, provider)))
	app.Get(
		"/api/oidc/callback",
		adaptor.HTTPHandler(rp.CodeExchangeHandler(rp.UserinfoCallback(callbackHandler), provider)),
	)
}

type OIDCAuthOutput struct {
	Status    int
	Location  string      `header:"Location"`
	SetCookie http.Cookie `header:"Set-Cookie"`
}

func S256CodeChallenge(code_verifier string) string {
	sum := sha256.Sum256([]byte(code_verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

type OIDCState struct {
	CodeVerifier string
	AuthTime     time.Time
}

// user ID to token state
var GitlabStates map[string]*OIDCState = map[string]*OIDCState{}

func gitlabRedirectURI() string {
	return fmt.Sprintf("%s/api/gitlab/callback", CONFIG.Listen.BaseURL)
}

type OIDCConfig struct {
	Issuer                string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint,omitempty"`
}

func (c OIDCConfig) Validate(issuer string) error {
	// TODO: check that issuer supports everything we need
	// for now we’re just assuming
	if issuer != c.Issuer {
		getLogger().Printf("wrong issuer in response from oidc well-known: %s != %s", issuer, c.Issuer)
		return errors.New("wrong issuer value in response from oidc well-known endpoint")
	}
	return nil
}

func NewOIDCProvider(issuer string, client_id string, scopes string, redirect_uri string, savefunc func(ctx context.Context, token OAuthTokenResponse) error) *OIDCProvider {

	provider := &OIDCProvider{
		ClientID:      client_id,
		Scopes:        scopes,
		RedirectURI:   redirect_uri,
		CodeVerifiers: map[string]*OIDCState{},
		SaveToken:     savefunc,
	}

	go provider.CleanupStates()

	return provider
}

func (provider *OIDCProvider) EnsureConfig(ctx context.Context) error {
	if provider.Config != nil {
		return nil
	}

	well_known_url := provider.Issuer + "/.well-known/openid-configuration"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, well_known_url, bytes.NewBufferString(""))
	if err != nil {
		return err
	}
	resp, err := Do(http.DefaultClient, req)
	if err != nil {
		return err
	}
	defer Closer(resp.Body)

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&provider.Config)
	if err != nil {
		return err
	}

	err = provider.Config.Validate(provider.Issuer)
	return err
}

var GITLAB_PROVIDER *OIDCProvider = nil

func GitlabProvider() *OIDCProvider {
	if GITLAB_PROVIDER == nil {
		scopes := strings.Join(CONFIG.Mugen.Gitlab.Scopes, "+")
		GITLAB_PROVIDER = NewOIDCProvider(CONFIG.Mugen.Gitlab.Server, CONFIG.Mugen.Gitlab.ClientID, scopes, gitlabRedirectURI(), saveGitlabToken)
	}
	return GITLAB_PROVIDER
}

type OIDCProvider struct {
	ClientID      string
	Issuer        string
	RedirectURI   string
	Config        *OIDCConfig
	Scopes        string
	CodeVerifiers map[string]*OIDCState // state → code verifier
	SaveToken     func(ctx context.Context, token OAuthTokenResponse) error
}

func (provider OIDCProvider) AuthorizeURI() string {
	return provider.Config.AuthorizationEndpoint
}

func (provider OIDCProvider) TokenURI() string {
	return provider.Config.TokenEndpoint
}

func (provider OIDCProvider) UserinfoURI() string {
	return provider.Config.UserinfoEndpoint
}

// delete unused states after some time
func (provider *OIDCProvider) CleanupStates() {
	for {
		time.Sleep(60 * time.Second)

		delete_before := time.Now().Add(-300 * time.Second)
		for k, v := range provider.CodeVerifiers {
			if v.AuthTime.Before(delete_before) {
				provider.CodeVerifiers[k] = nil
			}
		}
	}
}

func (provider *OIDCProvider) Auth(ctx context.Context, _ *struct{}) (*OIDCAuthOutput, error) {
	if !CONFIG.Mugen.Gitlab.IsSetup() {
		return nil, errors.New("gitlab client is not set up")
	}

	err := provider.EnsureConfig(ctx)
	if err != nil {
		return nil, err
	}

	state_verifier := authState()
	state_challenge := S256CodeChallenge(state_verifier)

	code_verifier := uuid.New().String()
	code_challenge := S256CodeChallenge(code_verifier)

	loc := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&state=%s&scope=%s&code_challenge=%s&code_challenge_method=S256",
		provider.AuthorizeURI(),
		provider.ClientID,
		provider.RedirectURI,
		state_challenge,
		provider.Scopes,
		code_challenge,
	)

	provider.CodeVerifiers[state_verifier] = &OIDCState{
		CodeVerifier: code_verifier,
		AuthTime:     time.Now(),
	}

	return &OIDCAuthOutput{
		Status:   http.StatusTemporaryRedirect,
		Location: loc,
		SetCookie: http.Cookie{
			Name:  "oidc_state",
			Value: state_verifier,
		},
	}, nil
}

type OIDCAuthCallbackInput struct {
	State         string `query:"state"`
	Code          string `query:"code"`
	StateVerifier string `cookie:"oidc_state"`
}

type OIDCAuthCallbackOutput struct {
	Status   int
	Location string `header:"Location"`
}

func (provider *OIDCProvider) getTokenCode(ctx context.Context, code_verifier string, code string, token_data *OAuthTokenResponse) error {
	url := fmt.Sprintf(
		"%s?client_id=%s&code=%s&grant_type=authorization_code&redirect_uri=%s&code_verifier=%s",
		provider.TokenURI(),
		provider.ClientID,
		code,
		gitlabRedirectURI(),
		code_verifier,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(""))
	if err != nil {
		return err
	}
	resp, err := Do(http.DefaultClient, req)
	if err != nil {
		return err
	}
	defer Closer(resp.Body)

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(token_data)
	return err
}

func (provider *OIDCProvider) Callback(ctx context.Context, input *OIDCAuthCallbackInput) (*OIDCAuthCallbackOutput, error) {
	err := provider.EnsureConfig(ctx)
	if err != nil {
		return nil, err
	}
	if S256CodeChallenge(input.StateVerifier) != input.State {
		return nil, huma.Error400BadRequest("Bad Request")
	}
	oidc_state := provider.CodeVerifiers[input.StateVerifier]
	if oidc_state == nil {
		return nil, huma.Error500InternalServerError("unknown state")
	}
	code_verifier := oidc_state.CodeVerifier
	provider.CodeVerifiers[input.StateVerifier] = nil

	token_data := OAuthTokenResponse{}
	err := provider.getTokenCode(ctx, code_verifier, input.Code, &token_data)
	if err != nil {
		return nil, err
	}

	err = provider.SaveToken(ctx, token_data)
	if err != nil {
		return nil, err
	}

	return &OIDCAuthCallbackOutput{
		Status:   http.StatusTemporaryRedirect,
		Location: "/",
	}, nil
}

func saveGitlabToken(ctx context.Context, token_data OAuthTokenResponse) error {
	token := OAuthToken{}
	err := setGitlabToken(GetDB(ctx), token_data, &token)
	if err != nil {
		return err
	}

	err = initOlderKarasExports(ctx)
	return err
}

func setGitlabToken(db *gorm.DB, token_data OAuthTokenResponse, token *OAuthToken) error {
	token.Server = CONFIG.Mugen.Gitlab.Server
	token.ClientID = CONFIG.Mugen.Gitlab.ClientID
	token.AccessToken = token_data.AccessToken
	token.RefreshToken = token_data.RefreshToken
	token.ExpiresAt = time.Now().Add(time.Duration(token_data.ExpiresIn) * time.Second)

	getLogger().Printf("new token expires at %d\n", token.ExpiresAt.Unix())

	err := db.Save(token).Error
	return DBErrToHumaErr(err)
}

func getGitlabToken(db *gorm.DB, token *OAuthToken) error {
	err := db.Where(&OAuthToken{
		Server:   CONFIG.Mugen.Gitlab.Server,
		ClientID: CONFIG.Mugen.Gitlab.ClientID,
	}).First(token).Error
	if err != nil {
		return err
	}

	if time.Now().Add(time.Duration(60) * time.Second).After(token.ExpiresAt) {
		return refreshGitlabToken(db.Statement.Context, db, token)
	}

	return err
}

func setDummyExports(db *gorm.DB) error {
	var karas []KaraInfoDB
	err := db.Scopes(CurrentKaras).Where(
		"id NOT IN (?) AND id NOT IN (?)",
		db.Table("mugen_exports").Select("kara_id AS id"),
		db.Table("mugen_imports").Select("kara_id AS id"),
	).Find(&karas).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	if len(karas) == 0 {
		return nil
	}

	kara_exports := []MugenExport{}
	for _, kara := range karas {
		mugen_export := MugenExport{KaraID: kara.ID, GitlabIssue: -1}
		kara_exports = append(kara_exports, mugen_export)
	}
	err = db.Create(&kara_exports).Error
	if err != nil {
		return err
	}

	return nil
}

func initOlderKarasExports(ctx context.Context) error {
	db := GetDB(ctx)

	var exportedKara MugenExport
	err := db.Where("gitlab_issue > 0").First(exportedKara).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// set a dummy export for older karas so they don’t get reexported
			// assuming that it is already done
			return setDummyExports(db)
		} else {
			return err
		}
	}

	// if we already have exported karas then we should catch up to the latest ones
	return exportRemainingKaras(ctx, db)
}

type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int    `json:"created_at"`
}

func refreshGitlabToken(ctx context.Context, db *gorm.DB, token *OAuthToken) error {
	getLogger().Printf("refreshing %s token\n", CONFIG.Mugen.Gitlab.Server)
	url := fmt.Sprintf(
		"%s/oauth/token?client_id=%s&refresh_token=%s&grant_type=refresh_token",
		CONFIG.Mugen.Gitlab.Server,
		CONFIG.Mugen.Gitlab.ClientID,
		token.RefreshToken,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(""))
	if err != nil {
		return err
	}
	resp, err := Do(http.DefaultClient, req)
	if err != nil {
		return err
	}
	defer Closer(resp.Body)

	if resp.StatusCode > 300 {
		buf := make([]byte, 2048)
		n, err := resp.Body.Read(buf)
		if err != nil {
			return err
		}

		getLogger().Printf("gitlab response: %+v\n%s", resp, buf[:n])
		return fmt.Errorf("gitlab responded with status code %d", resp.StatusCode)
	}

	token_data := &OAuthTokenResponse{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(token_data)
	if err != nil {
		return err
	}

	return setGitlabToken(db, *token_data, token)
}

func callbackHandler(
	w http.ResponseWriter,
	r *http.Request,
	tokens *oidc.Tokens[*oidc.IDTokenClaims],
	state string,
	rp rp.RelyingParty,
	info *oidc.UserInfo,
) {
	sub := info.Subject
	if CONFIG.OIDC.IDClaim != "" {
		sub = info.Claims[CONFIG.OIDC.IDClaim].(string)
	}

	expiresAt := time.Now().Add(time.Hour)
	_, signed, err := CreateJwtForUser(r.Context(), sub, &expiresAt, info)
	if err != nil {
		getLogger().Print(err)
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    sessionCookieName,
		Value:   signed,
		Path:    "/",
		Expires: expiresAt,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func CreateJwtForUser(
	ctx context.Context,
	sub string,
	expiresAt *time.Time,
	info *oidc.UserInfo,
) (*jwt.Token, string, error) {
	user, err := getOrCreateUser(ctx, sub, info)
	if err != nil {
		return nil, "", err
	}
	claims := KaraberusClaims{}
	if info != nil {
		claims.SetUserInfo(info)
	}
	claims.Subject = user.ID
	claims.IssuedAt = oidc.Time(time.Now().Unix())
	if expiresAt != nil {
		claims.Expiration = oidc.Time(expiresAt.Unix())
	}
	claims.Scopes = AllScopes
	claims.IsAdmin = user.Admin
	return createJwt(claims)
}

func createJwt(claims KaraberusClaims) (*jwt.Token, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(CONFIG.OIDC.JwtSignKey))
	return token, signed, err
}

func checkBasicAuth(token string) bool {
	if !CONFIG.Mugen.BasicAuth.isSetup() {
		return false
	}
	return token != CONFIG.Mugen.BasicAuth.Token()
}

func authError(ctx huma.Context, err error) {
	basicSecurity := false

	for _, opScheme := range ctx.Operation().Security {
		_, basicSecurityFound := opScheme["basic"]
		basicSecurity = basicSecurity || basicSecurityFound
	}

	if basicSecurity && CONFIG.Mugen.BasicAuth.isSetup() {
		ctx.SetHeader("WWW-Authenticate", "Basic realm=\"karaberus\"")
	}

	getLogger().Println(err)
	ctx.SetStatus(http.StatusUnauthorized)
}

func authMiddleware(ctx huma.Context, next func(huma.Context)) {
	var user *User = nil
	var scopes *Scopes = nil

	// Check for a token in the request.
	token, err := getRequestToken(ctx)
	if err != nil {
		authError(ctx, err)
		return
	}

	switch token.Type {
	// Bearer token
	case KaraberusBearerAuth:
		user, scopes, err = getUserScopesFromApiToken(ctx.Context(), token.Value)

		if err != nil {
			authError(ctx, err)
			return
		}

		ctx = huma.WithValue(ctx, currentUserCtxKey, *user)

	// Cookie/OIDC
	case KaraberusJWTAuth:
		user, scopes, err = getUserScopesFromJwt(ctx.Context(), token.Value)

		if err != nil {
			authError(ctx, err)
			return
		}

		ctx = huma.WithValue(ctx, currentUserCtxKey, *user)

	// Basic auth
	case KaraberusBasicAuth:
		if !checkBasicAuth(token.Value) {
			ctx.SetStatus(http.StatusForbidden)
			return
		}

		// no user
	}

	if checkOperationSecurity(ctx, user, scopes, token) {
		next(ctx)
	} else {
		ctx.SetStatus(http.StatusForbidden)
	}
}

type KaraberusAuthType string

var KaraberusBasicAuth KaraberusAuthType = "Basic"
var KaraberusBearerAuth KaraberusAuthType = "Bearer"
var KaraberusJWTAuth KaraberusAuthType = "Cookie"

type KaraberusAuthorization struct {
	Type  KaraberusAuthType
	Value string
}

var BASIC_AUTH_PREFIX string = "Basic "
var JWT_AUTH_PREFIX string = "JWT "

func getRequestToken(ctx huma.Context) (KaraberusAuthorization, error) {
	authHeader := ctx.Header("authorization")
	if authHeader != "" {
		if strings.HasPrefix(authHeader, oidc.BearerToken) {
			token := strings.TrimPrefix(authHeader, oidc.PrefixBearer)
			return KaraberusAuthorization{KaraberusBearerAuth, token}, nil
		}
		if strings.HasPrefix(authHeader, JWT_AUTH_PREFIX) {
			token := strings.TrimPrefix(authHeader, JWT_AUTH_PREFIX)
			return KaraberusAuthorization{KaraberusJWTAuth, token}, nil
		}
		if strings.HasPrefix(authHeader, BASIC_AUTH_PREFIX) {
			token := strings.TrimPrefix(authHeader, oidc.PrefixBearer)
			return KaraberusAuthorization{KaraberusBasicAuth, token}, nil
		}
		return KaraberusAuthorization{}, errors.New("invalid authorization header")
	} else {
		cookie, err := huma.ReadCookie(ctx, sessionCookieName)
		if err != nil {
			return KaraberusAuthorization{}, err
		}
		return KaraberusAuthorization{KaraberusJWTAuth, cookie.Value}, nil
	}
}

func getUserScopesFromApiToken(ctx context.Context, token string) (*User, *Scopes, error) {
	db := GetDB(ctx)
	apiToken := TokenV2{}
	if err := db.Preload(clause.Associations).Where(&TokenV2{Token: token}).First(&apiToken).Error; err != nil {
		return nil, nil, err
	}
	return &apiToken.User, &apiToken.Scopes, nil
}

func getUserScopesFromJwt(ctx context.Context, token string) (*User, *Scopes, error) {
	// claims := KaraberusClaims{} // FIXME: deserialization is not working with Scopes
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(CONFIG.OIDC.JwtSignKey), nil
	})
	if err != nil || !jwtToken.Valid {
		return nil, nil, &KaraberusError{"invalid token"}
	}

	sub, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return nil, nil, err
	}

	db := GetDB(ctx)
	user := User{ID: sub}
	if err := db.Preload(clause.Associations).First(&user).Error; err != nil {
		return nil, nil, err
	}

	// FIXME: replace with proper deserialization
	scopes := Scopes{}
	for key, val := range claims {
		if key == "kara" {
			scopes.Kara = val.(bool)
		}
		if key == "kara_ro" {
			scopes.KaraRO = val.(bool)
		}
		if key == "user" {
			scopes.User = val.(bool)
		}
	}

	return &user, &scopes, nil
}

func checkOperationSecurity(ctx huma.Context, user *User, scopes *Scopes, token KaraberusAuthorization) bool {
	oidcSecurity := false
	basicSecurity := false
	roles := []string{}
	opScopes := []string{}

	for _, opScheme := range ctx.Operation().Security {
		_, oidcSecurityFound := opScheme["oidc"]
		oidcSecurity = oidcSecurity || oidcSecurityFound

		_, basicSecurityFound := opScheme["basic"]
		basicSecurity = basicSecurity || basicSecurityFound

		opScopes = append(opScopes, opScheme["scopes"]...)
		roles = append(roles, opScheme["roles"]...)
	}

	adminRoute := slices.Contains(roles, "admin")

	if adminRoute && (user == nil || !user.Admin) {
		return false
	}

	// public endpoints
	if !oidcSecurity && !basicSecurity && opScopes == nil {
		return true
	}

	switch token.Type {
	case KaraberusBasicAuth:
		if basicSecurity {
			return true
		}
	case KaraberusJWTAuth:
		if oidcSecurity {
			return true
		}
	case KaraberusBearerAuth:
		for _, v := range opScopes {
			if scopes == nil || !scopes.HasScope(v) {
				return false
			}
		}

		// we have all the scopes to proceed
		return true
	}

	return false
}
