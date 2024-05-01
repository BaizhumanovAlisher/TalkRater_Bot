package helpers

import (
	"TalkRater_Bot/internal/data"
	"bytes"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env string `yaml:"env" env-required:"true"`

	ClearDbForNewConference bool             `yaml:"clear_db" env-default:"true"`
	ConferenceConfig        ConferenceConfig `yaml:"conference" env-required:"true"`
	SecretPath              SecretPath       `yaml:"secret"` // no parsing in config file is required
	DatabaseConfig          DatabaseConfig   `yaml:"database" env-required:"true"`

	TgBotSettings TgBotSettings    `yaml:"tg_bot_settings"`
	Conference    *data.Conference `yaml:"-"`
}

func MustLoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH_TG_BOT")

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can not read config: %s", err)
	}

	dbPassword := LoadOneSecret(cfg.SecretPath.DatabasePasswordPathFile)
	cfg.DatabaseConfig.compile(dbPassword)

	tokenUser := LoadOneSecret(cfg.SecretPath.TgTokenUserPathFile)
	tokenAdminPanel := LoadOneSecret(cfg.SecretPath.TgTokenAdminPanelPathFile)

	cfg.TgBotSettings.TokenUser = tokenUser
	cfg.TgBotSettings.TokenAdminPanel = tokenAdminPanel

	return &cfg
}

type ConferenceConfig struct {
	Name              string `yaml:"name" env-required:"true"`
	URL               string `yaml:"url" env-required:"true"`
	StartTime         string `yaml:"start_time" env-required:"true"`
	EndTime           string `yaml:"end_time" env-required:"true"`
	EndEvaluationTime string `yaml:"end_evaluation_time" env-required:"true"`
}

type SecretPath struct {
	DatabasePasswordPathFile  string `env:"DB_PASSWORD_FILE" env-required:"true"`
	TgTokenUserPathFile       string `env:"TG_API_TOKEN_USER_FILE" env-required:"true"`
	TgTokenAdminPanelPathFile string `env:"TG_API_TOKEN_ADMIN_FILE" env-required:"true"`
}

func LoadOneSecret(pathToFile string) string {
	fileBytes, err := os.ReadFile(pathToFile)
	if err != nil {
		log.Fatalf("can not read secret: %s", err)
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
