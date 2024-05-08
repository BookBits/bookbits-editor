package helpers

import (
	"github.com/BookBits/bookbits-editor/internal/handlers"
	"github.com/BookBits/bookbits-editor/internal/middlewares"
	"github.com/gofiber/fiber/v3"
)

func SetupHandlers(app *fiber.App) {
	app.Static("/static", "./public/")
	
	app.Get("/test", handlers.TestPage)
	app.Post("/test/increment", handlers.TestIncrement)
	app.Get("/check_db", handlers.CheckDB)

	app.Get("/", handlers.IndexHandler, middlewares.AuthMiddleware)
	app.Post("/refresh", handlers.RefreshSession)
	app.Get("/login", handlers.LoginPageHandler)
	app.Post("/login", handlers.Login)

	app.Get("/app", handlers.AppHomeHandler, middlewares.AuthMiddleware, middlewares.AuthOnlyRoute)
}
