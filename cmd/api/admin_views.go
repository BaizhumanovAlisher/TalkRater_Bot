package main

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"os"
	"strings"
	"talk_rater_bot/internal/templates"
	"talk_rater_bot/internal/templates/admin"
	"time"
)

const opStartAndHelpAdmin = "admin_views.startAndHelpAdmin"

func (app *application) startAndHelpAdmin(c tele.Context) error {
	app.logger.Info(opStartAndHelpAdmin,
		slog.String("username", c.Sender().Username))

	return c.Send(app.adminTemplates.Render(admin.StartInfo, nil))
}

const opSubmit = "admin_views.submitSchedule"

func (app *application) submitSchedule(c tele.Context) error {
	file := c.Message().Document
	if file == nil {
		app.logger.Error(opSubmit,
			slog.String("username", c.Sender().Username),
			slog.String("error", "file does not exist"),
		)

		return c.Send(app.adminTemplates.Render(admin.SubmitError, &templates.TemplateData{Error: "required to exist document in input"}))
	}

	if !strings.HasSuffix(file.FileName, ".csv") {
		app.logger.Info(opSubmit,
			slog.String("username", c.Sender().Username),
			slog.String("file name", file.FileName),
			slog.String("info", "name file should end `.csv`"),
		)

		return c.Send(app.adminTemplates.Render(admin.SubmitError, &templates.TemplateData{Error: "name file should end `.csv`"}))
	}

	filePath := fmt.Sprintf("%s%s%s_%s", app.pathTmp, string(os.PathSeparator), time.Now().Format("2006-01-02_15-04-05"), file.FileName)
	err := app.adminBot.Download(file.MediaFile(), filePath)
	if err != nil {
		app.logger.Error(opSubmit,
			slog.String("username", c.Sender().Username),
			slog.String("file path", filePath),
			slog.String("error", err.Error()),
		)

		return c.Send(app.adminTemplates.Render(admin.SubmitError, &templates.TemplateData{Error: "can not save file"}))
	}

	return c.Send(app.adminTemplates.Render(admin.SubmitSuccess, nil))
}
