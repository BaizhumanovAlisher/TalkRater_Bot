package view

import (
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"talk_rater_bot/internal/data"
	"talk_rater_bot/internal/helpers"
	"talk_rater_bot/internal/templates"
)

const opViewConference = "info.viewConference"

func (app *Application) viewConference(c tele.Context) error {
	app.Logger.Info(opViewConference, slog.String("username", c.Sender().Username))

	return c.Send(app.Templates.Render(templates.ConferenceTmpl,
		&templates.TemplateData{Conference: convertConf(app.Conference, app.TimeParser)}))
}

func convertConf(conf *data.Conference, parser *helpers.TimeParser) *templates.Conference {
	return &templates.Conference{
		Name:      conf.Name,
		URL:       conf.URL,
		StartTime: parser.ConvertTime(conf.StartTime),
		EndTime:   parser.ConvertTime(conf.EndTime),
	}
}
