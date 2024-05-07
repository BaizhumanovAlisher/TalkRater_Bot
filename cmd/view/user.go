package view

import (
	tele "gopkg.in/telebot.v3"
	"log/slog"
)

func (app *Application) helloWorld(c tele.Context) error {
	app.Logger.Info("request",
		slog.Int("args", len(c.Args())),
		slog.Int64("id", c.Sender().ID))
	return c.Send("Hello!")
}

func (app *Application) viewConference(context tele.Context) error {
	//todo:
	panic("implement me")
}
