package view

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log/slog"
	"talk_rater_bot/cmd/controllers"
	"talk_rater_bot/internal/databases"
	"talk_rater_bot/internal/helpers"
	"talk_rater_bot/internal/templates"
)

type Application struct {
	Logger     *slog.Logger
	UserBot    *tele.Bot
	AdminBot   *tele.Bot
	AdminDB    *databases.AdminDB
	Templates  *templates.Templates
	TimeParser *helpers.TimeParser
	PathTmp    string
	Controller *controllers.Controller
}

func (app *Application) Run() {
	go app.UserBot.Start()
	go app.AdminBot.Start()
}

func (app *Application) sendError(c tele.Context, err any) error {
	str := fmt.Sprintf("%+v", err)

	if len(str) > 1000 {
		return c.Send(app.Templates.Render(templates.Error, &templates.TemplateData{Error: str[:1000] + "...\nСообщение слишком длинное"}))
	}

	return c.Send(app.Templates.Render(templates.Error, &templates.TemplateData{Error: str}))
}
