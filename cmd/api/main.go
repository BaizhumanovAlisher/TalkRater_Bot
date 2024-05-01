package main

import (
	tele "gopkg.in/telebot.v3"
	"log"
	"log/slog"
	"talk_rater_bot/internal/helpers"
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

	user := &tele.User{ID: 12345678} // Replace with the actual user ID

	// Define the notification message
	notification := "This is your notification message."

	// Send the notification
	_, err = b.Send(user, notification)
}
