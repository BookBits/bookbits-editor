package helpers

import (
	"github.com/BookBits/bookbits-editor/internal/handlers"
	"github.com/BookBits/bookbits-editor/internal/middlewares"
	"github.com/gofiber/fiber/v3"
)

func SetupHandlers(app *fiber.App) {
	app.Static("/", "./public/")
	
	app.Get("/test", handlers.TestPage)
	app.Post("/test/increment", handlers.TestIncrement)
	app.Get("/check_db", handlers.CheckDB)

	app.Get("/", handlers.IndexPage, middlewares.AuthMiddleware)
	app.Get("/login", handlers.LoginPage)
	app.Post("/login", handlers.Login)
}
