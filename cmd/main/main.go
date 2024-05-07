package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	adminController "talk_rater_bot/cmd/controllers/admin"
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
	logger.Info(op, slog.String("info", "Start building application..."))

	userBot, err := helpers.SetupBot(cfg.TgBotSettings.TokenUser, cfg.TgBotSettings.Timeout)
	checkError(err, logger)

	adminBot, err := helpers.SetupBot(cfg.TgBotSettings.TokenAdminPanel, cfg.TgBotSettings.Timeout)
	checkError(err, logger)

	templatesMap, err := templates.NewTemplates(cfg.EnvVars.TemplatePath, templates.FilesName)
	checkError(err, logger)

	adminDB := databases.NewAdminDB(cfg.TgBotSettings.Admins)
	db, err := gorm.Open(postgres.Open(cfg.DatabaseConfig.CompiledFullPath),
		&gorm.Config{Logger: helpers.NewSlogLoggerDB(logger)})
	checkError(err, logger)

	dbHelper := databases.NewPrepareDBHelper(db,
		cfg.Conference, cfg.CleanupDBForNewConference, &cfg.DatabaseConfig, cfg.EnvVars.PathTmp)
	err = dbHelper.PrepareDB()
	checkError(err, logger)

	adminContr := adminController.NewController(db, cfg.TimeParser, cfg.Conference)

	app := view.Application{
		Logger:          logger,
		UserBot:         userBot,
		AdminBot:        adminBot,
		AdminDB:         adminDB,
		Templates:       templatesMap,
		TimeParser:      cfg.TimeParser,
		PathTmp:         cfg.EnvVars.PathTmp,
		AdminController: adminContr,
		Conference:      cfg.Conference,
	}

	app.Routes()
	app.Run()

	logger.Info(op, slog.String("info", "Start application..."))
	logger.Info(op, slog.String("info", "conference"),
		slog.String("name", cfg.Conference.Name),
		slog.String("url", cfg.Conference.URL),
		slog.Time("start", cfg.Conference.StartTime),
		slog.Time("end", cfg.Conference.EndTime),
		slog.Time("end evaluation", cfg.Conference.EndEvaluationTime),
	)

	_ = helpers.WaitSignal()

	logger.Info(op, slog.String("info", "Stop application. Started to stop"))

	app.UserBot.Stop()
	app.AdminBot.Stop()
	sqlDB, err := db.DB()
	if err != nil {
		logger.Warn(op,
			slog.String("info", "problem to close database connections"),
			slog.String("error", err.Error()),
		)
	} else {
		sqlDB.Close()
	}

	logger.Info(op, slog.String("info", "Stop application. Ended to stop"))
}

func checkError(err error, logger *slog.Logger) {
	if err != nil {
		logger.Warn(op, slog.String("error", err.Error()))
		panic(err)
	}
}
