package view

import (
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"talk_rater_bot/internal/templates"
	"time"
)

func (app *Application) recoverPanic(next tele.HandlerFunc) tele.HandlerFunc {
	const op = "middleware.recoverPanic"

	return func(c tele.Context) error {
		defer func() {
			if r := recover(); r != nil {
				app.Logger.Warn(op,
					slog.String("username", c.Sender().Username),
					slog.Any("panic", r))

				_ = c.Send("internal server error")
			}
		}()

		return next(c)
	}
}

func (app *Application) measureTime(next tele.HandlerFunc) tele.HandlerFunc {
	const op = "middleware.measureTime"

	return func(c tele.Context) error {
		timeStart := time.Now()

		result := next(c)

		duration := time.Now().Sub(timeStart)

		app.Logger.Info(op,
			slog.String("username", c.Sender().Username),
			slog.Duration("duration", duration),
		)

		return result
	}
}

func (app *Application) checkAdmin(next tele.HandlerFunc) tele.HandlerFunc {
	const op = "middleware.checkAdmin"

	return func(c tele.Context) error {
		username := c.Sender().Username
		if app.AdminDB.IsAdmin(username) {
			return next(c)
		} else {
			app.Logger.Info(op,
				slog.String("username", username),
				slog.String("info", "failed authorization"))
			return c.Send(app.Templates.Render(templates.AccessDeniedError, nil))
		}
	}
}

func (app *Application) checkUser(next tele.HandlerFunc) tele.HandlerFunc {
	const op = "middleware.checkUser"

	return func(c tele.Context) error {
		exists, err := app.Controller.UserExists(c.Sender().ID)
		if err != nil {
			app.Logger.Error(op,
				slog.String("username", c.Sender().Username),
				slog.String("error", err.Error()))
			return c.Send("проблема с авторизацией")
		}

		if !exists {
			app.Logger.Info(op,
				slog.String("username", c.Sender().Username),
				slog.String("info", "failed authorization"))

			return c.Send(app.Templates.Render(templates.UserAuthorization, nil))
		}

		return next(c)
	}
}
