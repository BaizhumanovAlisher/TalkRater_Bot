package view

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

func (app *Application) startAndHelpAdmin(c tele.Context) error {
	app.Logger.Info(opStartAndHelpAdmin,
		slog.String("username", c.Sender().Username))

	return c.Send(app.AdminTemplates.Render(admin.StartInfo, nil))
}

const opSubmit = "admin_views.submitSchedule"

func (app *Application) submitSchedule(c tele.Context) error {
	file := c.Message().Document
	if file == nil {
		app.Logger.Error(opSubmit,
			slog.String("username", c.Sender().Username),
			slog.String("error", "file does not exist"),
		)

		return c.Send(app.AdminTemplates.Render(admin.SubmitError, &templates.TemplateData{Error: "в сообщении должен быть файл"}))
	}

	if !strings.HasSuffix(file.FileName, ".csv") {
		app.Logger.Info(opSubmit,
			slog.String("username", c.Sender().Username),
			slog.String("file name", file.FileName),
			slog.String("info", "name file should end `.csv`"),
		)

		return c.Send(app.AdminTemplates.Render(admin.SubmitError, &templates.TemplateData{Error: "имя файла должно заканчивать на `.csv`"}))
	}

	filePath := fmt.Sprintf("%s%s%s_%s", app.PathTmp, string(os.PathSeparator), time.Now().Format("2006-01-02_15-04-05"), file.FileName)
	err := app.AdminBot.Download(file.MediaFile(), filePath)
	if err != nil {
		app.Logger.Error(opSubmit,
			slog.String("username", c.Sender().Username),
			slog.String("file path", filePath),
			slog.String("error", err.Error()),
		)

		return c.Send(app.AdminTemplates.Render(admin.SubmitError, &templates.TemplateData{Error: "не смог сохранить файл"}))
	}

	return c.Send(app.AdminTemplates.Render(admin.SubmitSuccess, nil))
}
