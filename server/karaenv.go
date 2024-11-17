package server

import (
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type KaraberusListenConfig struct {
	Host      string `envkey:"HOST" default:"127.0.0.1"`
	Port      int    `envkey:"PORT" default:"8888"`
	BaseURL   string `envkey:"BASE_URL"`
	Profiling bool   `envkey:"PROFILING"`
}

func (c KaraberusListenConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type KaraberusOIDCConfig struct {
	Issuer       string   `envkey:"ISSUER"`
	ClientID     string   `envkey:"CLIENT_ID"`
	ClientSecret string   `envkey:"CLIENT_SECRET"`
	Scopes       []string `envkey:"SCOPES" separator:" " default:"openid profile email groups"`
	IDClaim      string   `envkey:"ID_CLAIM"`
	GroupsClaim  string   `envkey:"GROUPS_CLAIM" default:"groups"`
	AdminGroup   string   `envkey:"ADMIN_GROUP" default:"admin"`
	JwtSignKey   string   `envkey:"JWT_SIGN_KEY"`
}

type KaraberusMugenGitlabConfig struct {
	Server       string   `envkey:"SERVER" default:"https://gitlab.com"`
	ProjectID    string   `envkey:"PROJECT_ID"`
	ClientID     string   `envkey:"CLIENT_ID"`
	ClientSecret string   `envkey:"CLIENT_SECRET"`
	Scopes       []string `envkey:"SCOPES" separator:" " default:"api"`
	ImportTag    string   `envkey:"IMPORT_TAG" default:"Import Japan7"`
	IssueLabels  []string `envkey:"LABELS" separator:"," default:"To Add"`
}

func (conf KaraberusMugenGitlabConfig) IsSetup() bool {
	return conf.ProjectID != "" && conf.ClientID != "" && conf.ClientSecret != ""
}

type BasicAuthConfig struct {
	Username string `envkey:"USERNAME"`
	Password string `envkey:"PASSWORD"`
}

func (basic BasicAuthConfig) isSetup() bool {
	return basic.Username != "" && basic.Password != ""
}

func (basic BasicAuthConfig) Token() string {
	return base64.RawStdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", basic.Username, basic.Password)),
	)
}

type KaraberusMugenConfig struct {
	Gitlab    KaraberusMugenGitlabConfig `env_prefix:"GITLAB"`
	BasicAuth BasicAuthConfig            `env_prefix:"BASIC"`
}

func (c KaraberusOIDCConfig) Validate() error {
	if c.Issuer == "" {
		return &KaraberusError{"OIDC issuer is not set"}
	}
	if c.ClientID == "" {
		return &KaraberusError{"OIDC client ID is not set"}
	}
	if c.ClientSecret == "" {
		return &KaraberusError{"OIDC client secret is not set"}
	}
	if c.JwtSignKey == "" {
		return &KaraberusError{"JWT signing key is not set"}
	}
	if c.GroupsClaim == "" {
		return &KaraberusError{"Groups claim is not set"}
	}
	if c.AdminGroup == "" {
		return &KaraberusError{"Admin group is not set"}
	}

	return nil
}

type KaraberusS3Config struct {
	Endpoints  []string `envkey:"ENDPOINT" separator:" "`
	KeyID      string   `envkey:"KEYID"`
	Secret     string   `envkey:"SECRET"`
	Secure     bool     `envkey:"SECURE"`
	BucketName string   `envkey:"BUCKET_NAME" default:"karaberus"`
}

type KaraberusDBConfig struct {
	Driver string `envkey:"DRIVER" default:"sqlite"`
	DSN    string `envkey:"DSN" default:"user=karaberus password=karaberus dbname=karaberus port=5123 sslmode=disable TimeZone=UTC"`
	File   string `envkey:"FILE" default:"karaberus.db"`
}

type KaraberusDakaraConfig struct {
	BaseURL string `envkey:"BASE_URL"`
	Token   string `envkey:"TOKEN"`
}

type KaraberusConfig struct {
	S3        KaraberusS3Config     `env_prefix:"S3"`
	OIDC      KaraberusOIDCConfig   `env_prefix:"OIDC"`
	Listen    KaraberusListenConfig `env_prefix:"LISTEN"`
	DB        KaraberusDBConfig     `env_prefix:"DB"`
	Dakara    KaraberusDakaraConfig `env_prefix:"DAKARA"`
	Mugen     KaraberusMugenConfig  `env_prefix:"MUGEN"`
	UIDistDir string                `envkey:"UI_DIST_DIR" default:"/usr/share/karaberus/ui_dist"`
	Webhooks  []string              `envkey:"WEBHOOKS" separator:" " example:"discord=<url1> discord=<url2> json=<url3>"`
}

func getEnvDefault(name string, defaultValue string) string {
	envVar := os.Getenv(name)
	if envVar != "" {
		return envVar
	}

	return defaultValue
}

func getFieldValue(field_type reflect.StructField, prefix string) string {
	envkey := field_type.Tag.Get("envkey")
	if envkey == "" {
		panic(fmt.Sprintf("envkey is not set for field %s", field_type.Name))
	}
	default_value := field_type.Tag.Get("default")
	return getEnvDefault(prefix+envkey, default_value)
}

func setConfigValue(config_value reflect.Value, config_type reflect.Type, prefix string) {
	for i := 0; i < config_type.NumField(); i++ {
		field_type := config_type.Field(i)
		field := config_value.FieldByName(field_type.Name)

		switch field_type.Type {
		case reflect.TypeOf([]string{}):
			value := getFieldValue(field_type, prefix)
			sep := field_type.Tag.Get("separator")
			arrval := strings.Split(value, sep)
			field.Set(reflect.ValueOf(arrval))
		case reflect.TypeOf(""):
			field.SetString(getFieldValue(field_type, prefix))
		case reflect.TypeOf(0):
			value := getFieldValue(field_type, prefix)
			int_value, err := strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
			field.SetInt(int64(int_value))
		case reflect.TypeOf(true):
			value := getFieldValue(field_type, prefix)
			field.SetBool(value != "" && strings.ToLower(value) != "false" && value != "0")
		default:
			if field_type.Type.Kind() == reflect.Struct {
				field_prefix := prefix + field_type.Tag.Get("env_prefix") + "_"
				setConfigValue(field, field_type.Type, field_prefix)
			} else {
				panic(fmt.Sprintf("unknown field type for field %s: %+v", field_type.Name, field_type.Type.Kind()))
			}
		}
	}
}

func getKaraberusConfig() KaraberusConfig {
	config := KaraberusConfig{}

	config_value := reflect.ValueOf(&config).Elem()
	config_type := reflect.TypeOf(config)

	setConfigValue(config_value, config_type, "KARABERUS_")

	if config.Listen.BaseURL == "" {
		// default to listen address
		config.Listen.BaseURL = "http://" + config.Listen.Addr()
		getLogger().Printf("Base URL implicitly set to %s", config.Listen.BaseURL)
	}

	return config
}

var CONFIG = getKaraberusConfig()
