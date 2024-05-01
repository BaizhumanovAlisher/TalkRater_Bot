package main

import (
	"TalkRater_Bot/internal/helpers"
	tele "gopkg.in/telebot.v3"
	"log"
	"log/slog"
)

func main() {
	cfg := helpers.MustLoadConfig()

	logger := helpers.SetupLogger(cfg.Env)
	slog.SetDefault(logger)
	logger.Info("Start application...")

	logger.Info("current conference",
		slog.String("name", cfg.Conference.Name),
		slog.String("url", cfg.Conference.URL),
		slog.Time("start", cfg.Conference.StartTime),
		slog.Time("end", cfg.Conference.EndTime),
		slog.Time("end evaluation", cfg.Conference.EndEvaluationTime),
	)

	pref := tele.Settings{
		Token:  cfg.TgBotSettings.TokenUser,
		Poller: &tele.LongPoller{Timeout: cfg.TgBotSettings.Timeout},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Start()
}
