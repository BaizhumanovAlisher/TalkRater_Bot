package config

import (
	"TalkRater_Bot/internal/data"
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	Env string `yaml:"env" env-required:"true"`

	ClearDbForNewConference bool              `yaml:"clear_db" env-default:"true"`
	ConferenceConfig        *ConferenceConfig `yaml:"conference" env-required:"true"`
	SecretPath              *SecretPath       `yaml:"-"`
	DatabaseConfig          *DatabaseConfig   `yaml:"database" env-required:"true"`

	TgBotSettings *TgBotSettings   `yaml:"tg_bot_settings"`
	Conference    *data.Conference `yaml:"-"`
}

func MustLoadConfig() *Config {

	return nil
}

type ConferenceConfig struct {
	Name              string `yaml:"name" env-required:"true"`
	URL               string `yaml:"url" env-required:"true"`
	StartTime         string `yaml:"start_time" env-required:"true"`
	EndTime           string `yaml:"end_time" env-required:"true"`
	EndEvaluationTime string `yaml:"end_evaluation_time" env-required:"true"`
}

type SecretPath struct {
	DatabasePasswordPathFile        string `env:"DB_PASSWORD_FILE" env-required:"true"`
	TelegramTokenUserPathFile       string `env:"TG_API_TOKEN_USER_FILE" env-required:"true"`
	TelegramTokenAdminPanelPathFile string `env:"TG_API_TOKEN_ADMIN_FILE" env-required:"true"`
}

func LoadOneSecret(pathToFile string) string {
	fileBytes, err := os.ReadFile(os.Getenv(pathToFile))
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes.TrimSpace(fileBytes))
}

type DatabaseConfig struct {
	Host         string `yaml:"host" env-required:"true"`
	Port         int    `yaml:"port" env-required:"true"`
	User         string `yaml:"user" env-required:"true"`
	DatabaseName string `yaml:"database_name" env-required:"true"`
	// Password will be downloaded from secrets

	CompiledFullPath string
}

func (dc *DatabaseConfig) compile(password string) {
	dc.CompiledFullPath = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dc.Host,
		dc.Port,
		dc.User,
		password,
		dc.DatabaseName,
	)
}

type TgBotSettings struct {
	Admins          []string      `yaml:"admins" env-required:"true"`
	Timeout         time.Duration `yaml:"timeout" env-required:"true"`
	TokenUser       string        `yaml:"-"`
	TokenAdminPanel string        `yaml:"-"`
}
