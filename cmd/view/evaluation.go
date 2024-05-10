package view

import (
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"strings"
	"talk_rater_bot/internal/templates"
)

const (
	opEvaluationZero       = "evaluation.evaluationZero"
	evaluationZeroPrefix   = "evaluation_0|"
	evaluationFirstPrefix  = "evaluation_1|"
	evaluationSecondPrefix = "evaluation_2|"
)

func (app *Application) evaluationZero(c tele.Context) error {
	data, _ := strings.CutPrefix(c.Callback().Data, "\f"+evaluationZeroPrefix)

	selector := &tele.ReplyMarkup{}
	selector.Inline(selector.Row(
		selector.Data("1", evaluationFirstPrefix, data, "1"),
		selector.Data("2", evaluationFirstPrefix, data, "2"),
		selector.Data("3", evaluationFirstPrefix, data, "3"),
	), selector.Row(
		selector.Data("4", evaluationFirstPrefix, data, "4"),
		selector.Data("5", evaluationFirstPrefix, data, "5"),
	))

	app.Logger.Info(opEvaluationZero, slog.String("username", c.Sender().Username))

	return c.Send(app.Templates.Render(templates.EvaluationZero, nil), selector)
}

const opEvaluationFirst = "evaluation.evaluationFirst"

func (app *Application) evaluationFirst(c tele.Context) error {
	data, _ := strings.CutPrefix(c.Callback().Data, "\f"+evaluationFirstPrefix)

	selector := &tele.ReplyMarkup{}
	selector.Inline(selector.Row(
		selector.Data("1", evaluationSecondPrefix, data, "1"),
		selector.Data("2", evaluationSecondPrefix, data, "2"),
		selector.Data("3", evaluationSecondPrefix, data, "3"),
	), selector.Row(
		selector.Data("4", evaluationSecondPrefix, data, "4"),
		selector.Data("5", evaluationSecondPrefix, data, "5"),
	))

	app.Logger.Info(opEvaluationFirst, slog.String("username", c.Sender().Username))

	return c.Send(app.Templates.Render(templates.EvaluationFirst, nil), selector)
}
