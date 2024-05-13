package view

import (
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"talk_rater_bot/internal/templates"
)

func (app *Application) helloUser() tele.HandlerFunc {
	const op = "user.helloUser"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		log := log.With(slog.String("username", c.Sender().Username))
		log.Info("")

		return c.Send(app.Templates.Render(templates.StartInfoUser, nil))
	}

}

func (app *Application) identicalInfo(c tele.Context) error {
	panic("not implemented")
}
