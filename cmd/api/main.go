package main

import (
	tele "gopkg.in/telebot.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"talk_rater_bot/internal/databases"
	"talk_rater_bot/internal/helpers"
	"talk_rater_bot/internal/templates"
	"talk_rater_bot/internal/templates/admin"
	"talk_rater_bot/internal/templates/user"
)

func main() {
	cfg := helpers.MustLoadConfig()

	logger := helpers.SetupLogger(cfg.Env, cfg.EnvVars.PathLogs)
	slog.SetDefault(logger)

	userBot := setupBot(cfg.TgBotSettings.TokenUser, cfg.TgBotSettings.Timeout)
	adminBot := setupBot(cfg.TgBotSettings.TokenAdminPanel, cfg.TgBotSettings.Timeout)

	adminTemplates, err := templates.NewTemplates(cfg.EnvVars.TemplatePath, admin.DirectoryName, admin.FilesName)
	if err != nil {
		log.Fatal(err)
	}
	userTemplates, err := templates.NewTemplates(cfg.EnvVars.TemplatePath, user.DirectoryName, user.FilesName)
	if err != nil {
		log.Fatal(err)
	}

	adminDB := databases.NewAdminDB(cfg.TgBotSettings.Admins)
	db, err := gorm.Open(postgres.Open(cfg.DatabaseConfig.CompiledFullPath), &gorm.Config{
		Logger: helpers.NewSlogLoggerDB(logger),
	})

	err = databases.AutoMigrateAllModels(db)
	if err != nil {
		log.Fatal(err)
	}

	app := application{
		logger:         logger,
		userBot:        userBot,
		adminBot:       adminBot,
		adminDB:        adminDB,
		adminTemplates: adminTemplates,
		userTemplates:  userTemplates,
	}

	app.routes()
	app.run()

	logger.Info("Start application...")
	logger.Info("conference",
		slog.String("name", cfg.Conference.Name),
		slog.String("url", cfg.Conference.URL),
		slog.Time("start", cfg.Conference.StartTime),
		slog.Time("end", cfg.Conference.EndTime),
		slog.Time("end evaluation", cfg.Conference.EndEvaluationTime),
	)

	_ = helpers.WaitSignal()

	logger.Info("Stop application...")
	app.userBot.Stop()
	app.adminBot.Stop()
}

type application struct {
	logger         *slog.Logger
	userBot        *tele.Bot
	adminBot       *tele.Bot
	adminDB        *databases.AdminDB
	adminTemplates *templates.Templates
	userTemplates  *templates.Templates
}

func (app *application) run() {
	go app.userBot.Start()
	go app.adminBot.Start()
}
