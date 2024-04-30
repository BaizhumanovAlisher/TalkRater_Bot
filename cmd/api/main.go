package main

import (
	"TalkRater_Bot/internal/config"
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	secrets := config.MustLoadSecret()

	pref := tele.Settings{
		Token:  string(secrets.TelegramTokenUser),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
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
