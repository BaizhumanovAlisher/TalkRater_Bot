package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	adminController "talk_rater_bot/cmd/controllers"
	"talk_rater_bot/cmd/view"
	"talk_rater_bot/internal/databases"
	"talk_rater_bot/internal/helpers"
	"talk_rater_bot/internal/templates"
)

const op = "main.main"

func main() {
	cfg := helpers.MustLoadConfig()

	logger := helpers.SetupLogger(cfg.Env, cfg.EnvVars.PathTmp)
	slog.SetDefault(logger)

	log := logger.With(slog.String("op", op))
	log.Info("Start building application...")

	userBot, err := helpers.SetupBot(cfg.TgBotSettings.TokenUser, cfg.TgBotSettings.Timeout)
	checkError(err, log)

	adminBot, err := helpers.SetupBot(cfg.TgBotSettings.TokenAdminPanel, cfg.TgBotSettings.Timeout)
	checkError(err, log)

	templatesMap, err := templates.NewTemplates(cfg.EnvVars.TemplatePath, templates.FilesName)
	checkError(err, log)

	adminDB := databases.NewAdminDB(cfg.TgBotSettings.Admins)
	db, err := gorm.Open(postgres.Open(cfg.DatabaseConfig.CompiledFullPath),
		&gorm.Config{Logger: helpers.NewSlogLoggerDB(logger)})
	checkError(err, log)

	dbHelper := databases.NewPrepareDBHelper(db, cfg.Conference,
		cfg.CleanupDBForNewConference, &cfg.DatabaseConfig, cfg.EnvVars.PathTmp)
	err = dbHelper.PrepareDB()
	checkError(err, log)

	adminContr := adminController.NewController(db, cfg.TimeParser, cfg.Conference)

	app := view.Application{
		Logger:     logger,
		UserBot:    userBot,
		AdminBot:   adminBot,
		AdminDB:    adminDB,
		Templates:  templatesMap,
		TimeParser: cfg.TimeParser,
		PathTmp:    cfg.EnvVars.PathTmp,
		Controller: adminContr,
	}

	app.Routes()
	app.Run()

	log.Info("Start application...")
	log.Info("conference", slog.String("name", cfg.Conference.Name),
		slog.String("url", cfg.Conference.URL),
		slog.Time("start", cfg.Conference.StartTime), slog.Time("end", cfg.Conference.EndTime),
		slog.Time("end evaluation", cfg.Conference.EndEvaluationTime),
	)

	_ = helpers.WaitSignal()

	log.Info("Stop application. Started to stop")

	app.UserBot.Stop()
	app.AdminBot.Stop()
	sqlDB, err := db.DB()
	if err != nil {
		log.Warn(err.Error())
	} else {
		sqlDB.Close()
	}

	log.Info("Stop application. Ended to stop")
}

func checkError(err error, logger *slog.Logger) {
	if err != nil {
		logger.Warn(err.Error())
		panic(err)
	}
}
