package main

import (
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"talk_rater_bot/internal/templates/admin"
)

const opStartAndHelpAdmin = "admin_bot.startAndHelpAdmin"

func (app *application) startAndHelpAdmin(c tele.Context) error {
	app.logger.Info(opStartAndHelpAdmin,
		slog.String("username", c.Sender().Username))

	return c.Send(
		app.adminTemplates.Render(admin.StartInfo))
}

const opSubmit = "admin_bot.submitSchedule"

func (app *application) submitSchedule(c tele.Context) error {
	app.logger.Info(opSubmit,
		slog.String("username", c.Sender().Username))

	return c.Send(
		app.adminTemplates.Render(admin.SubmitInfo))
}
