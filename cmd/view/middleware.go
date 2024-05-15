package view

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"talk_rater_bot/internal/data"
	"talk_rater_bot/internal/templates"
	"time"
)

type middlewareController interface {
	UserExists(id int64) (bool, error)
	RetrieveSession(chatID int64) (*data.Session, bool, error)
}

func (app *Application) recoverPanic(next tele.HandlerFunc) tele.HandlerFunc {
	const op = "middleware.recoverPanic"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		defer func() {
			if r := recover(); r != nil {
				log.Warn(fmt.Sprintf("%v", r), slog.String("username", c.Sender().Username))

				_ = c.Send("internal server error")
			}
		}()

		return next(c)
	}
}

func (app *Application) measureTime(next tele.HandlerFunc) tele.HandlerFunc {
	const op = "middleware.measureTime"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		log := log.With(slog.String("username", c.Sender().Username))
		timeStart := time.Now()

		result := next(c)

		duration := time.Now().Sub(timeStart)

		log.Info(duration.String())

		return result
	}
}

func (app *Application) checkAdmin(next tele.HandlerFunc) tele.HandlerFunc {
	const op = "middleware.checkAdmin"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		username := c.Sender().Username
		if app.AdminDB.IsAdmin(username) {
			return next(c)
		} else {
			log.Info("failed authorization", slog.String("username", username))
			return c.Send(app.Templates.Render(templates.AccessDeniedError, nil))
		}
	}
}

func (app *Application) checkUser(next tele.HandlerFunc) tele.HandlerFunc {
	const op = "middleware.checkUser"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		log := log.With(slog.String("username", c.Sender().Username))
		exists, err := app.MiddlewareController.UserExists(c.Sender().ID)

		if err != nil {
			log.Error(err.Error())
			return c.Send("проблема с авторизацией")
		}

		if !exists {
			log.Info("failed authorization")
			return c.Send(app.Templates.Render(templates.UserAuthorization, nil))
		}

		return next(c)
	}
}
