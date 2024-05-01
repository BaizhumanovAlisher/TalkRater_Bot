package main

import (
	"TalkRater_Bot/internal/helpers"
	tele "gopkg.in/telebot.v3"
	"log"
	"log/slog"
)

func main() {
	config := helpers.MustLoadConfig()

	logger := helpers.SetupLogger(config.Env)
	slog.SetDefault(logger)
	logger.Info("Start bot...")

	pref := tele.Settings{
		Token:  config.TgBotSettings.TokenUser,
		Poller: &tele.LongPoller{Timeout: config.TgBotSettings.Timeout},
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
