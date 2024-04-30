package main

import (
	configpackage "TalkRater_Bot/internal/config"
	tele "gopkg.in/telebot.v3"
	"log"
)

func main() {
	config := configpackage.MustLoadConfig()

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
