package main

import (
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"talk_rater_bot/internal/templates/admin"
)

const opHelloWorldHandler = "admin_bot.helloWorldAdmin"

func (app *application) helloWorldAdmin(c tele.Context) error {
	app.logger.Info(opHelloWorldHandler,
		slog.String("id", c.Sender().Username))

	return c.Send(
		app.adminTemplates.Render(admin.StartInfo))
}
