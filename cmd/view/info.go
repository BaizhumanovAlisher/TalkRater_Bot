package view

import (
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"strconv"
	"strings"
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

const (
	opViewLectures = "info.viewSchedule"
	limit          = 10
)

func (app *Application) viewSchedule(c tele.Context) error {
	maxCountPage, err := app.Controller.CountPageInLectures(limit)
	pageNumber := generatePageNumber(c.Callback(), maxCountPage)

	lectures, err := app.Controller.GetSchedule(limit, (pageNumber-1)*limit)
	if err != nil {
		app.Logger.Warn(opViewLectures, slog.String("error", err.Error()))
		return c.Send(app.Templates.Render(templates.SubmitError,
			&templates.TemplateData{Error: err.Error()}))
	}

	app.Logger.Info(opViewConference, slog.String("username", c.Sender().Username))

	return c.Send(app.generateResponse(lectures, pageNumber))
}

func (app *Application) generateResponse(lectures []*data.Lecture, pageNumber int) (string, *tele.ReplyMarkup) {
	schedule := make([]*templates.Lecture, len(lectures))
	selector := &tele.ReplyMarkup{}
	btnNext := selector.Data("➡", "next", strconv.Itoa(pageNumber))
	btnPrev := selector.Data("⬅", "prev", strconv.Itoa(pageNumber))
	buttons := make([]tele.Btn, len(lectures))

	for i, lecture := range lectures {
		number := i + 1 + (pageNumber-1)*limit

		schedule[i] = convertShortLecture(lecture, number, app.TimeParser)
		buttons[i] = selector.Data(strconv.Itoa(number), "id", strconv.FormatInt(lecture.ID, 10))
	}

	mid := len(buttons) / 2

	selector.Inline(
		selector.Row(buttons[:mid]...),
		selector.Row(buttons[mid:]...),
		selector.Row(btnPrev, btnNext),
	)

	return app.Templates.Render(templates.Schedule, &templates.TemplateData{Schedule: schedule}), selector
}

func convertShortLecture(lecture *data.Lecture, number int, parser *helpers.TimeParser) *templates.Lecture {
	return &templates.Lecture{
		Number:    strconv.Itoa(number),
		Name:      lecture.Title,
		StartTime: parser.ConvertTime(lecture.Start),
	}
}

func generatePageNumber(c *tele.Callback, maxCountPage int64) (pageNumber int) {
	if c == nil {
		return 1
	}

	var numStr string
	var add bool

	txt, _ := strings.CutPrefix(c.Data, "\f")

	if strings.HasPrefix(txt, "prev|") {
		add = false
		numStr, _ = strings.CutPrefix(txt, "prev|")
	} else if strings.HasPrefix(txt, "next|") {
		add = true
		numStr, _ = strings.CutPrefix(txt, "next|")
	} else {
		panic("invalid callback data: " + txt)
	}

	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		panic(err)
	}

	if add {
		num++
	} else {
		num--
	}

	if num < 1 {
		return int(maxCountPage)
	} else if num > maxCountPage {
		return 1
	}

	return int(num)
}
