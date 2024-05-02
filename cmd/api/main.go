package main

import (
	tele "gopkg.in/telebot.v3"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"talk_rater_bot/internal/helpers"
	"time"
)

func main() {
	cfg := helpers.MustLoadConfig()

	logger := helpers.SetupLogger(cfg.Env)
	slog.SetDefault(logger)

	logger.Info("conference",
		slog.String("name", cfg.Conference.Name),
		slog.String("url", cfg.Conference.URL),
		slog.Time("start", cfg.Conference.StartTime),
		slog.Time("end", cfg.Conference.EndTime),
		slog.Time("end evaluation", cfg.Conference.EndEvaluationTime),
	)

	userBot := setupBot(cfg.TgBotSettings.TokenUser, cfg.TgBotSettings.Timeout)
	adminBot := setupBot(cfg.TgBotSettings.TokenAdminPanel, cfg.TgBotSettings.Timeout)

	app := application{
		logger:   logger,
		userBot:  userBot,
		adminBot: adminBot,
	}

	app.routes()
	app.run()
	logger.Info("Start application...")

	_ = waitSignal()

	logger.Info("Stop application...")
	app.userBot.Stop()
	app.adminBot.Stop()
}

type application struct {
	logger   *slog.Logger
	userBot  *tele.Bot
	adminBot *tele.Bot
}

func (app *application) run() {
	go app.userBot.Start()
	go app.adminBot.Start()
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

func waitSignal() os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	return <-stop
}
