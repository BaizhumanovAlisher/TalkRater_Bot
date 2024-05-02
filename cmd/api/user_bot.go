package main

import (
	tele "gopkg.in/telebot.v3"
	"log/slog"
)

func (app *application) helloWorld(c tele.Context) error {
	app.logger.Info("request",
		slog.Int("args", len(c.Args())),
		slog.Int64("id", c.Sender().ID))
	return c.Send("Hello!")
}
