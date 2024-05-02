package main

func (app *application) routes() {
	app.adminBotRoutes()
	app.userBotRoutes()
}

func (app *application) userBotRoutes() {
	app.userBot.Handle("/hello", app.helloWorld)
}

func (app *application) adminBotRoutes() {

}
