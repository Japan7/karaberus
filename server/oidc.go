package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OIDCProviderDiscovery struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	Issuer                string `json:"issuer"`
	JWKSURI               string `json:"jwks_uri"`
	TokenEndpoint         string `json:"token_endpoint"`
	ClientID              string `json:"client_id"`
	UserInfoEndpoint      string `json:"userinfo_endpoint"`
}

type OIDCAuthEndpointOutput struct {
	Body OIDCProviderDiscovery
}

func getOIDCProviderDiscovery() (*OIDCProviderDiscovery, error) {
	oidc_config_url := fmt.Sprintf("%s/.well-known/openid-configuration", CONFIG.OIDC.Issuer)

	resp, err := http.Get(oidc_config_url)
	if err != nil {
		return nil, err
	}

	data := &OIDCProviderDiscovery{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(data)
	data.ClientID = CONFIG.OIDC.ClientID

	return data, err
}

func getOIDCDiscovery(ctx context.Context, input *struct{}) (*OIDCAuthEndpointOutput, error) {
	oidc_config, err := getOIDCProviderDiscovery()
	if err != nil {
		return nil, err
	}

	resp := &OIDCAuthEndpointOutput{}
	resp.Body = *oidc_config

	return resp, nil
}
