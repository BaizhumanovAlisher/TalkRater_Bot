package view

import (
	tele "gopkg.in/telebot.v3"
	"strings"
)

func (app *Application) Routes() {
	app.UserBot.Use(app.recoverPanic, app.measureTime)

	app.UserBot.Handle("/conference", app.viewConference)
	app.UserBot.Handle("/schedule", app.viewConference)
	app.UserBot.Handle(tele.OnCallback, func(c tele.Context) error {
		if strings.HasPrefix(c.Callback().Data, "prev|") || strings.HasPrefix(c.Callback().Data, "next|") {
			return app.viewSchedule(c)
		}
		return c.Send("unimplemented method")
	})

	app.AdminBot.Use(app.recoverPanic, app.measureTime, app.checkAdmin)

	app.AdminBot.Handle("/start", app.startAndHelpAdmin)
	app.AdminBot.Handle("/help", app.startAndHelpAdmin)
	app.AdminBot.Handle(tele.OnDocument, app.submitSchedule)
	app.AdminBot.Handle("/export", app.exportEvaluations)
	app.AdminBot.Handle("/conference", app.viewConference)
	app.AdminBot.Handle("/schedule", app.viewSchedule)
	app.AdminBot.Handle(tele.OnCallback, func(c tele.Context) error {
		txt, _ := strings.CutPrefix(c.Callback().Data, "\f")
		if strings.HasPrefix(txt, "prev|") || strings.HasPrefix(txt, "next|") {
			return app.viewSchedule(c)
		}
		return c.Send("unimplemented method")
	})
}
