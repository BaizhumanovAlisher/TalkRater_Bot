package view

import (
	tele "gopkg.in/telebot.v3"
)

func (app *Application) Routes() {
	app.adminBotRoutes()
	app.userBotRoutes()
}

func (app *Application) userBotRoutes() {
	app.UserBot.Use(app.recoverPanic, app.measureTime)

	app.UserBot.Handle("/start", app.helloWorld)
}

func (app *Application) adminBotRoutes() {
	app.AdminBot.Use(app.recoverPanic, app.measureTime, app.checkAdmin)

	app.AdminBot.Handle("/start", app.startAndHelpAdmin)
	app.AdminBot.Handle("/help", app.startAndHelpAdmin)
	app.AdminBot.Handle(tele.OnDocument, app.submitSchedule)
}
