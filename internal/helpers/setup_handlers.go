package helpers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/BookBits/bookbits-editor/internal/handlers"
)

func SetupHandlers(app *fiber.App) {
	app.Static("/", "./public/")
	
	app.Get("/test", handlers.TestPage)
	app.Post("/test/increment", handlers.TestIncrement)
	app.Get("/check_db", handlers.CheckDB)
}
