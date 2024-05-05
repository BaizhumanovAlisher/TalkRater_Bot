package main

import (
	tele "gopkg.in/telebot.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"talk_rater_bot/cmd/view"
	"talk_rater_bot/internal/databases"
	"talk_rater_bot/internal/helpers"
	"talk_rater_bot/internal/templates"
	"talk_rater_bot/internal/templates/admin"
	"talk_rater_bot/internal/templates/user"
	"time"
)

const op = "main.main"

func main() {
	cfg := helpers.MustLoadConfig()

	logger := helpers.SetupLogger(cfg.Env, cfg.EnvVars.PathTmp)
	slog.SetDefault(logger)
	logger.Info(op, slog.String("info", "Start building application..."))

	userBot := setupBot(cfg.TgBotSettings.TokenUser, cfg.TgBotSettings.Timeout)
	adminBot := setupBot(cfg.TgBotSettings.TokenAdminPanel, cfg.TgBotSettings.Timeout)

	adminTemplates, err := templates.NewTemplates(cfg.EnvVars.TemplatePath, admin.DirectoryName, admin.FilesName)
	if err != nil {
		logger.Warn(op, slog.String("error", err.Error()))
		panic(err)
	}
	userTemplates, err := templates.NewTemplates(cfg.EnvVars.TemplatePath, user.DirectoryName, user.FilesName)
	if err != nil {
		logger.Warn(op, slog.String("error", err.Error()))
		panic(err)
	}

	adminDB := databases.NewAdminDB(cfg.TgBotSettings.Admins)
	db, err := gorm.Open(postgres.Open(cfg.DatabaseConfig.CompiledFullPath), &gorm.Config{
		Logger: helpers.NewSlogLoggerDB(logger),
	})

	dbHelper := databases.NewPrepareDBHelper(db, cfg.Conference, cfg.CleanupDBForNewConference, &cfg.DatabaseConfig, cfg.EnvVars.PathTmp)
	err = dbHelper.PrepareDB()
	if err != nil {
		logger.Warn(op, slog.String("error", err.Error()))
		panic(err)
	}

	app := view.Application{
		Logger:         logger,
		UserBot:        userBot,
		AdminBot:       adminBot,
		AdminDB:        adminDB,
		AdminTemplates: adminTemplates,
		UserTemplates:  userTemplates,
		TimeParser:     cfg.TimeParser,
		PathTmp:        cfg.EnvVars.PathTmp,
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
	sqlDB, _ := db.DB()
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

func setupBot(token string, timeout time.Duration) *tele.Bot {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: timeout},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
