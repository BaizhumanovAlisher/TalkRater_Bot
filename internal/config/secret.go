package config

import (
	"bytes"
	"log"
	"os"
)

var (
	EnvDBPasswordFile     = "DB_PASSWORD_FILE"
	EnvBotTokenUser       = "TG_API_TOKEN_USER_FILE"
	EnvBotTokenAdminPanel = "TG_API_TOKEN_ADMIN_FILE"
)

type Secret struct {
	DatabasePassword        []byte
	TelegramTokenUser       []byte
	TelegramTokenAdminPanel []byte
}

func MustLoadSecret() *Secret {
	return &Secret{
		DatabasePassword:        LoadOneSecret(EnvDBPasswordFile),
		TelegramTokenUser:       LoadOneSecret(EnvBotTokenUser),
		TelegramTokenAdminPanel: LoadOneSecret(EnvBotTokenAdminPanel),
	}
}

func LoadOneSecret(env string) []byte {
	path := os.Getenv(env)
	if path == "" {
		log.Fatalf("empty environment variable for %s", env)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return bytes.TrimSpace(data)
}
