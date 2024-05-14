package view

import (
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"talk_rater_bot/internal/data"
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

func (app *Application) identicalInfo() tele.HandlerFunc {
	const op = "user.identicalInfo"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		log := log.With(slog.String("username", c.Sender().Username))

		session := &data.Session{
			ChatID: c.Chat().ID,
			UserID: c.Sender().ID,
			Form:   data.UserIdenticalInfoForm,
		}

		err := app.Controller.SaveSession(session)
		if err != nil {
			log.Error(err.Error())
			return app.sendError(c, "проблема с сохранением формы")
		}

		return c.Send(app.Templates.Render(templates.UserAuthForm, nil))
	}
}

func (app *Application) identicalInfoForm() tele.HandlerFunc {
	const op = "user.identicalInfoForm"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		log := log.With(slog.String("username", c.Sender().Username))
		txt := c.Message().Text
		if txt == "" {
			log.Info("message is empty")
			return app.sendError(c, "сообщение не может быть пустым")
		}

		user := &data.User{
			ID:           c.Sender().ID,
			UserName:     c.Sender().Username,
			IdentityInfo: txt,
		}

		err := app.Controller.SaveUser(user)
		if err != nil {
			log.Error(err.Error())
			return c.Send("ошибка сохранения данных")
		}

		return c.Send("данные успешно сохранены")
	}
}
