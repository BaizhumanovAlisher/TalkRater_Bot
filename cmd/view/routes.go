package view

import (
	tele "gopkg.in/telebot.v3"
)

func (app *Application) Routes() {
	app.UserBot.Use(app.recoverPanic, app.measureTime)

	app.UserBot.Handle("/conference", app.viewConference)
	app.AdminBot.Handle("/schedule", app.viewConference)

	app.AdminBot.Use(app.recoverPanic, app.measureTime, app.checkAdmin)

	app.AdminBot.Handle("/start", app.startAndHelpAdmin)
	app.AdminBot.Handle("/help", app.startAndHelpAdmin)
	app.AdminBot.Handle(tele.OnDocument, app.submitSchedule)
	app.AdminBot.Handle("/export", app.exportEvaluations)
	app.AdminBot.Handle("/conference", app.viewConference)
	app.AdminBot.Handle("/schedule", app.viewSchedule)
}
