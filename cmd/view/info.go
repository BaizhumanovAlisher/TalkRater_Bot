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

func (app *Application) viewConference() tele.HandlerFunc {
	op := "info.viewConference"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		log := log.With(slog.String("username", c.Sender().Username))
		log.Info("")

		return c.Send(app.Templates.Render(templates.ConferenceTmpl,
			&templates.TemplateData{Conference: convertConf(app.Controller.GetCurrentConference(), app.TimeParser)}))
	}
}

func convertConf(conf *data.Conference, parser *helpers.TimeParser) *templates.Conference {
	return &templates.Conference{
		Name:      conf.Name,
		URL:       conf.URL,
		StartTime: parser.ConvertTime(conf.StartTime),
		EndTime:   parser.ConvertTime(conf.EndTime),
	}
}

const limit = 10

func (app *Application) viewSchedule() tele.HandlerFunc {
	op := "info.viewSchedule"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		log := log.With(slog.String("username", c.Sender().Username))
		maxCountPage, err := app.Controller.CountPageInLectures(limit)
		pageNumber := generatePageNumber(c.Callback(), maxCountPage)

		lectures, err := app.Controller.GetSchedule(limit, (pageNumber-1)*limit)
		if err != nil {
			log.Warn(err.Error())

			return c.Send(app.Templates.Render(templates.Error, &templates.TemplateData{Error: err.Error()}))
		}

		log.Info("")

		return c.Send(app.generateResponse(lectures, pageNumber))
	}

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
		buttons[i] = selector.Data(strconv.Itoa(number), "id",
			strconv.FormatInt(lecture.ID, 10), strconv.FormatInt(int64(pageNumber), 10))
	}

	if len(buttons) > 1 {

		mid := len(buttons) / 2

		selector.Inline(
			selector.Row(buttons[:mid]...),
			selector.Row(buttons[mid:]...),
			selector.Row(btnPrev, btnNext),
		)
	} else {
		selector.Inline(
			selector.Row(buttons...),
			selector.Row(btnPrev, btnNext),
		)
	}

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

func (app *Application) viewLecture() tele.HandlerFunc {
	op := "info.viewLecture"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		log := log.With(slog.String("username", c.Sender().Username))

		numbers, _ := strings.CutPrefix(c.Callback().Data, "\fid|")
		nums := strings.Split(numbers, "|")
		if len(nums) != 2 {
			log.Warn("callback error - only two args")
			return c.Send("должно быть два аргумента")
		}

		lectureID, err := strconv.ParseInt(nums[0], 10, 64)
		if err != nil {
			log.Warn("id must be number")
			return c.Send("id должно быть числом")
		}

		pageNumber, err := strconv.ParseInt(nums[1], 10, 64)
		if err != nil {
			log.Warn("page number must be number")
			return c.Send("номер страницы должен быть числом")
		}

		lecture, err := app.Controller.GetLecture(lectureID)
		if err != nil {
			log.Error(err.Error())
			return c.Send(app.Templates.Render(templates.Error,
				&templates.TemplateData{Error: "проблемы с базой данных"}))
		}

		app.Logger.Info(op, slog.String("username", c.Sender().Username))
		selector := &tele.ReplyMarkup{}
		selector.Inline(selector.Row(
			selector.Data("Вернуться", "next", strconv.FormatInt(pageNumber-1, 10)),
			selector.Data("Оценить", evaluationZeroPrefix, strconv.FormatInt(lectureID, 10)),
		))

		return c.Send(app.Templates.Render(templates.LectureTmpl,
			&templates.TemplateData{Lecture: convertFullLecture(lecture, app.TimeParser)}), selector)
	}
}

func convertFullLecture(lecture *data.Lecture, timeParser *helpers.TimeParser) *templates.Lecture {
	return &templates.Lecture{
		Name:      lecture.Title,
		Speaker:   lecture.Speaker,
		URL:       lecture.URL,
		StartTime: timeParser.ConvertTime(lecture.Start),
		EndTime:   timeParser.ConvertTime(lecture.End),
	}
}
