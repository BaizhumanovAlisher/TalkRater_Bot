package helpers

import (
	tele "gopkg.in/telebot.v3"
	"time"
)

func SetupBot(token string, timeout time.Duration) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: timeout},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}
	return b, nil
}
