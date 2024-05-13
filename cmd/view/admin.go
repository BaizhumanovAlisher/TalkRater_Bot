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

func (app *Application) helloAdmin() tele.HandlerFunc {
	op := "admin.helloAdmin"
	log := app.Logger.With("op", op)

	return func(c tele.Context) error {
		log.Info("", slog.String("username", c.Sender().Username))
		return c.Send(app.Templates.Render(templates.StartInfoAdmin, nil))
	}
}

func (app *Application) submitSchedule() tele.HandlerFunc {
	const opSubmit = "admin.submitSchedule"
	log := app.Logger.With(slog.String("op", opSubmit))

	return func(c tele.Context) error {
		log := log.With(slog.String("username", c.Sender().Username))

		file := c.Message().Document

		if !strings.HasSuffix(file.FileName, ".csv") {
			log.Info("name file should end `.csv`")
			return c.Send(app.Templates.Render(templates.Error,
				&templates.TemplateData{Error: "имя файла должно заканчивать на `.csv`"}))
		}

		filePath := app.generateFilePath(file.FileName)
		err := app.AdminBot.Download(file.MediaFile(), filePath)
		if err != nil {
			log.Error(err.Error())
			return c.Send(app.Templates.Render(templates.Error, &templates.TemplateData{Error: "не смог сохранить файл"}))
		}

		defer func() {
			err = os.Remove(filePath)
			if err != nil {
				log.Warn(err.Error(), slog.String("file path", filePath))
			}
		}()

		err = app.Controller.GenerateSchedule(filePath)
		if err != nil {
			log.Error(err.Error())
			return submitError(c, app, err)
		}

		return c.Send(app.Templates.Render(templates.SubmitSuccess, nil))
	}

}

func (app *Application) exportEvaluations() tele.HandlerFunc {
	const op = "admin.exportEvaluations"
	log := app.Logger.With(slog.String("op", op))

	return func(c tele.Context) error {
		log := log.With(slog.String("username", c.Sender().Username))
		evaluations, err := app.Controller.ExportEvaluations()

		if err != nil {
			log.Error(err.Error())
			return submitError(c, app, err)
		}

		jsonData, err := json.Marshal(evaluations)
		if err != nil {
			log.Error(err.Error())
			return submitError(c, app, err)
		}

		fileName := "evaluations.json"
		filePath := app.generateFilePath(fileName)

		err = os.WriteFile(filePath, jsonData, 0666)
		if err != nil {
			log.Error(err.Error())
			return submitError(c, app, err)
		}

		defer func() {
			err = os.Remove(filePath)
			if err != nil {
				log.Warn(err.Error(), slog.String("file path", filePath))
			}
		}()

		fileTG := &tele.Document{File: tele.FromDisk(filePath), FileName: fileName}
		log.Info("export file was send")
		return c.Send(fileTG)
	}
}

func (app *Application) generateFilePath(fileName string) string {
	return fmt.Sprintf("%s%s%s_%s", app.PathTmp, string(os.PathSeparator), time.Now().Format("2006-01-02_15-04-05"), fileName)
}

func submitError(c tele.Context, app *Application, err error) error {
	if len(err.Error()) > 1000 {
		return c.Send(app.Templates.Render(templates.Error, &templates.TemplateData{Error: err.Error()[:1000] + "...\nСообщение слишком длинное"}))
	}
	return c.Send(app.Templates.Render(templates.Error, &templates.TemplateData{Error: err.Error()}))
}
