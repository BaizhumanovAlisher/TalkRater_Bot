package main

import (
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
)

func (app *application) routes() {
	app.adminBotRoutes()
	app.userBotRoutes()
}

func (app *application) userBotRoutes() {
	app.userBot.Use(app.recoverPanic, app.measureTime)

	app.userBot.Handle("/start", app.helloWorld)
}

func (app *application) adminBotRoutes() {
	app.adminBot.Use(app.recoverPanic, app.measureTime, app.checkAdmin)

	app.adminBot.Handle("/start", app.startAndHelpAdmin)
	app.adminBot.Handle("/help", app.startAndHelpAdmin)
	app.adminBot.Handle("/submit", app.submitSchedule)
}

func setupBot(token string, timeout time.Duration) *tele.Bot {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: timeout},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
