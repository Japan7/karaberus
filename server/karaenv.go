package server

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type KaraberusListenConfig struct {
	Host string `envkey:"HOST"`
	Port int    `envkey:"PORT" default:"8888"`
}

func (c KaraberusListenConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type KaraberusOIDCConfig struct {
	Issuer   string `envkey:"ISSUER"`
	KeyID    string `envkey:"KEY_ID"`
	Key      string `envkey:"KEY"`
	IDClaim  string `envkey:"ID_CLAIM"`
	ClientID string `envkey:"CLIENT_ID"`
}

func (c KaraberusOIDCConfig) Validate() error {
	if c.Issuer == "" {
		return &KaraberusError{"OIDC issuer is not set"}
	}
	if c.Key == "" {
		return &KaraberusError{"OIDC key is not set"}
	}
	if c.ClientID == "" {
		return &KaraberusError{"OIDC client ID is not set"}
	}

	return nil
}

type KaraberusS3Config struct {
	Endpoint   string `envkey:"ENDPOINT"`
	KeyID      string `envkey:"KEYID"`
	Secret     string `envkey:"SECRET"`
	Secure     bool   `envkey:"SECURE"`
	BucketName string `envkey:"BUCKET_NAME" default:"karaberus"`
}

type KaraberusDBConfig struct {
	Driver string `envkey:"DRIVER" default:"sqlite"`
	DSN    string `envkey:"DSN" default:"user=karaberus password=karaberus dbname=karaberus port=5123 sslmode=disable TimeZone=UTC"`
	File   string `envkey:"FILE" default:"karaberus.db"`
	Delete bool   `envkey:"DELETE"`
}

type KaraberusConfig struct {
	S3                 KaraberusS3Config     `env_prefix:"S3"`
	OIDC               KaraberusOIDCConfig   `env_prefix:"OIDC"`
	Listen             KaraberusListenConfig `env_prefix:"LISTEN"`
	GENERATED_TEST_DIR string                `envkey:"GENERATED_TEST_DIR"`
	TEST_DIR           string                `envkey:"TEST_DIR"`
	DB                 KaraberusDBConfig     `env_prefix:"DB"`
	UIDistDir          string                `envkey:"UI_DIST_DIR" default:"/usr/share/karaberus/ui_dist"`
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

		switch field_type.Type.Kind() {
		case reflect.String:
			field.SetString(getFieldValue(field_type, prefix))
		case reflect.Int:
			value := getFieldValue(field_type, prefix)
			int_value, err := strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
			field.SetInt(int64(int_value))
		case reflect.Bool:
			field.SetBool(getFieldValue(field_type, prefix) != "")
		case reflect.Struct:
			field_prefix := prefix + field_type.Tag.Get("env_prefix") + "_"
			setConfigValue(field, field_type.Type, field_prefix)
		default:
			panic(fmt.Sprintf("unknown field type for field %s: %s", field_type.Name, field_type.Type.Name()))
		}
	}
}

func getKaraberusConfig() KaraberusConfig {
	config := KaraberusConfig{}

	config_value := reflect.ValueOf(&config).Elem()
	config_type := reflect.TypeOf(config)

	setConfigValue(config_value, config_type, "KARABERUS_")

	return config
}

var CONFIG = getKaraberusConfig()
