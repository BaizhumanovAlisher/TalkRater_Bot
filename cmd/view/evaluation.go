package view

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"strconv"
	"strings"
	"talk_rater_bot/internal/data"
	"talk_rater_bot/internal/templates"
	"talk_rater_bot/internal/validators"
	"time"
)

const (
	evaluationZeroPrefix   = "evaluation_0"
	evaluationFirstPrefix  = "evaluation_1"
	evaluationSecondPrefix = "evaluation_2"
)

type evaluationController interface {
	SaveEvaluation(eval *data.Evaluation) error
	SaveSession(session *data.Session) error
	SaveComment(id int64, comment string) error
	UserExists(id int64) (bool, error)
	GetLecture(id int64) (*data.Lecture, error)
	GetCurrentConference() *data.Conference
}

func (app *Application) evaluationZero() tele.HandlerFunc {
	const op = "evaluation.evaluationZero"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		dataInput, _ := strings.CutPrefix(c.Callback().Data, "\f"+evaluationZeroPrefix+"|")

		selector := &tele.ReplyMarkup{}
		selector.Inline(selector.Row(
			selector.Data("1", evaluationFirstPrefix, dataInput, "1"),
			selector.Data("2", evaluationFirstPrefix, dataInput, "2"),
			selector.Data("3", evaluationFirstPrefix, dataInput, "3"),
		), selector.Row(
			selector.Data("4", evaluationFirstPrefix, dataInput, "4"),
			selector.Data("5", evaluationFirstPrefix, dataInput, "5"),
		))

		log.Info("", slog.String("username", c.Sender().Username))

		return c.Send(app.Templates.Render(templates.EvaluationZero, nil), selector)
	}
}

func (app *Application) evaluationFirst() tele.HandlerFunc {
	const op = "evaluation.evaluationFirst"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		dataInput, _ := strings.CutPrefix(c.Callback().Data, "\f"+evaluationFirstPrefix+"|")

		selector := &tele.ReplyMarkup{}
		selector.Inline(selector.Row(
			selector.Data("1", evaluationSecondPrefix, dataInput, "1"),
			selector.Data("2", evaluationSecondPrefix, dataInput, "2"),
			selector.Data("3", evaluationSecondPrefix, dataInput, "3"),
		), selector.Row(
			selector.Data("4", evaluationSecondPrefix, dataInput, "4"),
			selector.Data("5", evaluationSecondPrefix, dataInput, "5"),
		))

		log.Info("", slog.String("username", c.Sender().Username))

		return c.Send(app.Templates.Render(templates.EvaluationFirst, nil), selector)
	}
}

func (app *Application) evaluationSecond() tele.HandlerFunc {
	const op = "evaluation.evaluationSecond"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		log := log.With(slog.String("username", c.Sender().Username), slog.String("data", c.Callback().Data))
		dataInput, _ := strings.CutPrefix(c.Callback().Data, "\f"+evaluationSecondPrefix+"|")

		args := strings.Split(dataInput, "|")
		if len(args) != 3 {
			log.Warn("must be 3 arg")
			return app.sendError(c, "количество аргументов в callback должно быть 3")
		}

		lectureID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Error(err.Error())
			return app.sendError(c, "проблемы с id доклада")
		}

		content, err := strconv.Atoi(args[1])
		if err != nil {
			log.Error(err.Error())
			return app.sendError(c, "проблемы с конвертацией string в int")
		}

		performance, err := strconv.Atoi(args[2])
		if err != nil {
			log.Error(err.Error())
			return app.sendError(c, "проблемы с конвертацией string в int")
		}

		evaluation := &data.Evaluation{
			UserID:           c.Sender().ID,
			LectureID:        lectureID,
			TypeEvaluation:   data.Correct,
			ScoreContent:     int8(content),
			ScorePerformance: int8(performance),
		}

		v := validators.New()
		data.ValidateEvaluation(v, evaluation)
		if !v.Valid() {
			return app.sendError(c, v.Errors)
		}

		now := time.Now()
		end := app.EvaluationController.GetCurrentConference().EndEvaluationTime
		if now.After(end) {
			log.Info("end time evaluation")
			return c.Send(fmt.Sprintf("время приема оценок закончилось %s назад", now.Sub(end).Abs().String()))
		}

		lecture, err := app.EvaluationController.GetLecture(lectureID)
		if err != nil {
			log.Error(err.Error())
			return app.sendError(c, err)
		}

		start := lecture.Start
		if now.Before(start) {
			log.Info("user tried evaluate before lecture start")
			return c.Send(fmt.Sprintf("до начала лекции осталось %s", now.Sub(start).Abs().String()))
		}

		err = app.EvaluationController.SaveEvaluation(evaluation)
		if err != nil {
			log.Error(err.Error())
			return app.sendError(c, "ошибка при сохранении оценки")
		}

		app.Logger.Info("created evaluation")

		session := &data.Session{
			ChatID:       c.Chat().ID,
			UserID:       c.Sender().ID,
			Form:         data.CommentForm,
			EvaluationID: evaluation.ID,
		}

		err = app.EvaluationController.SaveSession(session)
		if err != nil {
			log.Error(err.Error())
			return app.sendError(c, "оценка сохранена, но не создана форма для комментария")
		}

		return c.Send(app.Templates.Render(templates.EvaluationSecond, nil))
	}
}

func (app *Application) submitComment() tele.HandlerFunc {
	const op = "evaluation.submitComment"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		log := log.With(slog.String("username", c.Sender().Username))

		txt := c.Message().Text

		if len([]rune(txt)) < 4 {
			log.Info("short comment", "comment", txt)
			return c.Send("комментарий был проигнорирован, так как длина меньше 4 символов.")
		}

		session, ok := c.Get(data.SessionKey).(*data.Session)
		if !ok || session == nil {
			log.Error("session type error")
			return app.sendError(c, "проблема с сессиями")
		}

		err := app.EvaluationController.SaveComment(session.EvaluationID, txt)
		if err != nil {
			log.Error(err.Error())
			return app.sendError(c, "ошибка при сохранении комментария")
		}

		err = c.Send(app.Templates.Render(templates.CommentSuccess, nil))
		if err != nil {
			return err
		}

		exists, err := app.EvaluationController.UserExists(c.Sender().ID)
		if err != nil {
			return err
		}

		if !exists {
			return c.Send(app.Templates.Render(templates.UserAuthorization, nil))
		}

		return nil
	}
}
