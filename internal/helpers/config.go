package helpers

import (
	"bytes"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"talk_rater_bot/internal/data"
	"talk_rater_bot/internal/validators"
	"time"
)

type Config struct {
	Env string `yaml:"env" env-required:"true"`

	CleanupDBForNewConference bool             `yaml:"clear_db" env-default:"true"`
	ConferenceConfig          ConferenceConfig `yaml:"conference" env-required:"true"`
	EnvVars                   EnvVars          `yaml:"secret"` // no parsing in config file is required. only for env
	DatabaseConfig            DatabaseConfig   `yaml:"database" env-required:"true"`

	TgBotSettings TgBotSettings    `yaml:"tg_bot_settings"`
	Conference    *data.Conference `yaml:"-"`
	TimeParser    *TimeParser      `yaml:"-"`
}

func MustLoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH_TG_BOT")

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can not read config: %s", err)
	}

	dbPassword := LoadOneSecret(cfg.EnvVars.DatabasePasswordPathFile)
	cfg.DatabaseConfig.compile(dbPassword)

	tokenUser := LoadOneSecret(cfg.EnvVars.TgTokenUserPathFile)
	tokenAdminPanel := LoadOneSecret(cfg.EnvVars.TgTokenAdminPanelPathFile)

	cfg.TgBotSettings.TokenUser = tokenUser
	cfg.TgBotSettings.TokenAdminPanel = tokenAdminPanel
	cfg.TgBotSettings.validateAdmins()

	var err error
	cfg.TimeParser, err = NewTimeParserInMoscow()
	if err != nil {
		log.Fatalf("can not create time parser: %s", err)
	}

	cfg.MustLoadConference()

	return &cfg
}

type ConferenceConfig struct {
	Name              string `yaml:"name" env-required:"true"`
	URL               string `yaml:"url" env-required:"true"`
	StartTime         string `yaml:"start_time" env-required:"true"`
	EndTime           string `yaml:"end_time" env-required:"true"`
	EndEvaluationTime string `yaml:"end_evaluation_time" env-required:"true"`
}

type EnvVars struct {
	DatabasePasswordPathFile  string `env:"DB_PASSWORD_FILE" env-required:"true"`
	TgTokenUserPathFile       string `env:"TG_API_TOKEN_USER_FILE" env-required:"true"`
	TgTokenAdminPanelPathFile string `env:"TG_API_TOKEN_ADMIN_FILE" env-required:"true"`
	PathTmp                   string `env:"PATH_TMP"`
	TemplatePath              string `env:"TEMPLATE_PATH" env-required:"true"`
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
	Password     string `yaml:"-"` // Password will be downloaded from secrets

	CompiledFullPath string
}

func (dc *DatabaseConfig) compile(password string) {
	dc.Password = password

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

func (tbs *TgBotSettings) validateAdmins() {
	v := validators.New()

	v.Check(len(tbs.Admins) > 0, "admins count", "Count of admins must be greater than 0")
	v.Check(validators.Unique(tbs.Admins), "admins unique", "admins must be unique")

	for _, admin := range tbs.Admins {
		v.Check(admin != "", "admins not empty", "admins must not be empty")
	}

	if !v.Valid() {
		log.Fatalf("non valid admins: %v", tbs.Admins)
	}
}

func (cfg *Config) convertConference(confStr *ConferenceConfig) (*data.Conference, error) {
	startTime, err := cfg.TimeParser.ParseTimeString(confStr.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := cfg.TimeParser.ParseTimeString(confStr.EndTime)
	if err != nil {
		return nil, err
	}

	endEvaluationTime, err := cfg.TimeParser.ParseTimeString(confStr.EndEvaluationTime)
	if err != nil {
		return nil, err
	}

	return &data.Conference{
		Name:              confStr.Name,
		URL:               confStr.URL,
		StartTime:         startTime,
		EndTime:           endTime,
		EndEvaluationTime: endEvaluationTime,
	}, nil
}

func (cfg *Config) MustLoadConference() {
	conference, err := cfg.convertConference(&cfg.ConferenceConfig)
	if err != nil {
		log.Fatalf("can not convert conference: %s", err)
	}

	v := validators.New()
	data.ValidateConference(v, conference)
	if !v.Valid() {
		log.Fatalf("can not validate conference: %v", v.Errors)
	}
	cfg.Conference = conference
}
