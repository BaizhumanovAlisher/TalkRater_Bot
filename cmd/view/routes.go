package view

import (
	tele "gopkg.in/telebot.v3"
	"strings"
)

func (app *Application) Routes() {
	app.UserBot.Use(app.recoverPanic, app.measureTime)

	app.UserBot.Handle("/start", app.helloUser)
	app.UserBot.Handle("/help", app.helloUser)
	app.UserBot.Handle("/conference", app.viewConference)
	app.UserBot.Handle("/schedule", app.viewConference)
	app.UserBot.Handle(tele.OnCallback, app.callbackRouter)

	app.AdminBot.Use(app.recoverPanic, app.measureTime, app.checkAdmin)

	app.AdminBot.Handle("/start", app.helloAdmin)
	app.AdminBot.Handle("/help", app.helloAdmin)
	app.AdminBot.Handle(tele.OnDocument, app.submitSchedule)
	app.AdminBot.Handle("/export", app.exportEvaluations)
}

func (app *Application) callbackRouter(c tele.Context) error {
	txt, _ := strings.CutPrefix(c.Callback().Data, "\f")
	if strings.HasPrefix(txt, "prev|") || strings.HasPrefix(txt, "next|") {
		return app.viewSchedule(c)
	}
	if strings.HasPrefix(txt, "id|") {
		return app.viewLecture(c)
	}
	if strings.HasPrefix(txt, evaluationZeroPrefix) {
		return app.evaluationZero(c)
	}
	if strings.HasPrefix(txt, evaluationFirstPrefix) {
		return app.evaluationFirst(c)
	}

	app.Logger.Warn("unimplemented method")
	return c.Send("unimplemented method")
}
