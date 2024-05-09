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
		&templates.TemplateData{Conference: convertConf(app.Controller.GetCurrentConference(), app.TimeParser)}))
}

func convertConf(conf *data.Conference, parser *helpers.TimeParser) *templates.Conference {
	return &templates.Conference{
		Name:      conf.Name,
		URL:       conf.URL,
		StartTime: parser.ConvertTime(conf.StartTime),
		EndTime:   parser.ConvertTime(conf.EndTime),
	}
}

const opViewLectures = "info.viewSchedule"

func (app *Application) viewSchedule(c tele.Context) error {
	lectures, err := app.Controller.GetSchedule()
	if err != nil {
		app.Logger.Warn(opViewLectures, slog.String("error", err.Error()))
		return c.Send(app.Templates.Render(templates.SubmitError,
			&templates.TemplateData{Error: err.Error()}))
	}

	app.Logger.Info(opViewConference, slog.String("username", c.Sender().Username))

	return c.Send(app.Templates.Render(templates.Schedule,
		&templates.TemplateData{Schedule: convertShortSchedule(lectures, app.TimeParser)}))
}

func convertShortSchedule(lectures []*data.Lecture, parser *helpers.TimeParser) []*templates.Lecture {
	shortLectures := make([]*templates.Lecture, len(lectures))

	for i := 0; i < len(lectures); i++ {
		shortLectures[i] = &templates.Lecture{
			Name:      lectures[i].Title,
			StartTime: parser.ConvertTime(lectures[i].Start),
		}
	}

	return shortLectures
}
