package main

func (app *application) routes() {
	app.adminBotRoutes()
	app.userBotRoutes()
}

func (app *application) userBotRoutes() {
	app.userBot.Use(app.recoverPanic, app.measureTime)

	app.userBot.Handle("/start", app.helloWorld)
}

func (app *application) adminBotRoutes() {
	app.adminBot.Use(app.recoverPanic, app.measureTime, app.checkAdmin)

	app.adminBot.Handle("/start", app.helloWorldAdmin)
}
