package server

import (
	"fmt"
	"os"
	"strconv"
)

type KaraberusListenConfig struct {
	Host string
	Port int
}

func (c KaraberusListenConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type KaraberusOIDCConfig struct {
	Issuer   string
	KeyID    string
	Key      string
	IDClaim  string
	ClientID string
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
	Endpoint   string
	KeyID      string
	Secret     string
	Secure     bool
	BucketName string
}

type KaraberusDBConfig struct {
	File   string
	Delete bool
}

type KaraberusConfig struct {
	S3                 KaraberusS3Config
	OIDC               KaraberusOIDCConfig
	Listen             KaraberusListenConfig
	GENERATED_TEST_DIR string
	TEST_DIR           string
	DB                 KaraberusDBConfig
	UIDistDir          string
}

func getEnvDefault(name string, defaultValue string) string {
	envVar := os.Getenv("KARABERUS_" + name)
	if envVar != "" {
		return envVar
	}

	return defaultValue
}

func getKaraberusConfig() KaraberusConfig {
	config := KaraberusConfig{}

	config.Listen.Host = getEnvDefault("LISTEN_HOST", "")
	port, err := strconv.Atoi(getEnvDefault("LISTEN_PORT", "8888"))
	if err != nil {
		panic(err)
	}
	config.Listen.Port = port

	config.S3.Endpoint = getEnvDefault("S3_ENDPOINT", "")
	config.S3.KeyID = getEnvDefault("S3_KEYID", "")
	config.S3.Secret = getEnvDefault("S3_SECRET", "")
	config.S3.Secure = getEnvDefault("S3_SECURE", "") != ""
	config.S3.BucketName = getEnvDefault("BUCKET_NAME", "karaberus")

	config.OIDC.Issuer = getEnvDefault("OIDC_ISSUER", "")
	config.OIDC.KeyID = getEnvDefault("OIDC_KEY_ID", "")
	config.OIDC.Key = getEnvDefault("OIDC_KEY", "")
	config.OIDC.ClientID = getEnvDefault("OIDC_CLIENT_ID", "")
	config.OIDC.IDClaim = getEnvDefault("OIDC_ID_CLAIM", "")

	config.GENERATED_TEST_DIR = getEnvDefault("GENERATED_TEST_DIR", ".")
	config.TEST_DIR = getEnvDefault("TEST_DIR", ".")
	config.DB.File = getEnvDefault("DB_FILE", "karaberus.db")
	config.DB.Delete = getEnvDefault("DELETE_DB", "") != ""
	config.UIDistDir = getEnvDefault("UI_DIST_DIR", "/usr/share/karaberus/ui_dist")

	return config
}

var CONFIG = getKaraberusConfig()
