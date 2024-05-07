package view

import (
	"encoding/json"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"os"
	"strings"
	"talk_rater_bot/internal/templates"
	"time"
)

const opStartAndHelpAdmin = "admin.startAndHelpAdmin"

func (app *Application) startAndHelpAdmin(c tele.Context) error {
	app.Logger.Info(opStartAndHelpAdmin,
		slog.String("username", c.Sender().Username))

	return c.Send(app.Templates.Render(templates.StartInfoAdmin, nil))
}

const opSubmit = "admin.submitSchedule"

func (app *Application) submitSchedule(c tele.Context) error {
	file := c.Message().Document
	if file == nil {
		app.Logger.Error(opSubmit,
			slog.String("username", c.Sender().Username),
			slog.String("error", "file does not exist"),
		)

		return c.Send(app.Templates.Render(templates.SubmitError, &templates.TemplateData{Error: "в сообщении должен быть файл"}))
	}

	if !strings.HasSuffix(file.FileName, ".csv") {
		app.Logger.Info(opSubmit,
			slog.String("username", c.Sender().Username),
			slog.String("file name", file.FileName),
			slog.String("info", "name file should end `.csv`"),
		)

		return c.Send(app.Templates.Render(templates.SubmitError, &templates.TemplateData{Error: "имя файла должно заканчивать на `.csv`"}))
	}

	filePath := app.generateFilePath(file.FileName)
	err := app.AdminBot.Download(file.MediaFile(), filePath)
	if err != nil {
		app.Logger.Error(opSubmit,
			slog.String("username", c.Sender().Username),
			slog.String("file path", filePath),
			slog.String("error", err.Error()),
		)

		return c.Send(app.Templates.Render(templates.SubmitError, &templates.TemplateData{Error: "не смог сохранить файл"}))
	}

	err = app.AdminController.GenerateSchedule(filePath)
	if err != nil {
		app.Logger.Error(opSubmit,
			slog.String("username", c.Sender().Username),
			slog.String("file path", filePath),
			slog.String("error", err.Error()),
		)

		return c.Send(app.Templates.Render(templates.SubmitError, &templates.TemplateData{Error: err.Error()}))
	}

	err = os.Remove(filePath)
	if err != nil {
		app.Logger.Error(opSubmit,
			slog.String("username", c.Sender().Username),
			slog.String("file path", filePath),
			slog.String("error", err.Error()),
			slog.String("info", "problem to delete file"),
		)
	}

	return c.Send(app.Templates.Render(templates.SubmitSuccess, nil))
}

const opExport = "admin.exportEvaluations"

func (app *Application) exportEvaluations(c tele.Context) error {
	evaluations, err := app.AdminController.ExportEvaluations()

	if err != nil {
		app.Logger.Error(opExport,
			slog.String("username", c.Sender().Username),
			slog.String("error", err.Error()),
		)

		return c.Send(app.Templates.Render(templates.SubmitError, &templates.TemplateData{Error: err.Error()}))
	}

	jsonData, err := json.Marshal(evaluations)
	if err != nil {
		app.Logger.Error(opExport,
			slog.String("username", c.Sender().Username),
			slog.String("error", err.Error()),
		)

		return c.Send(app.Templates.Render(templates.SubmitError, &templates.TemplateData{Error: err.Error()}))
	}

	fileName := "evaluations.json"
	filePath := app.generateFilePath(fileName)

	err = os.WriteFile(filePath, jsonData, 0666)
	if err != nil {
		app.Logger.Error(opExport,
			slog.String("username", c.Sender().Username),
			slog.String("error", err.Error()),
		)

		return c.Send(app.Templates.Render(templates.SubmitError, &templates.TemplateData{Error: err.Error()}))
	}

	fileTG := &tele.Document{File: tele.FromDisk(filePath), FileName: fileName}
	err = c.Send(fileTG)
	app.Logger.Info(opExport,
		slog.String("username", c.Sender().Username),
		slog.String("info", "export file was send"),
	)

	err = os.Remove(filePath)
	if err != nil {
		app.Logger.Error(opExport,
			slog.String("username", c.Sender().Username),
			slog.String("file path", filePath),
			slog.String("error", err.Error()),
			slog.String("info", "problem to delete file"),
		)
	}

	return err
}

func (app *Application) generateFilePath(fileName string) string {
	return fmt.Sprintf("%s%s%s_%s", app.PathTmp, string(os.PathSeparator), time.Now().Format("2006-01-02_15-04-05"), fileName)
}
