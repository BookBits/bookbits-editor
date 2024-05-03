package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/BookBits/bookbits-editor/internal/helpers"
)

func main() {
	app := fiber.New()
	helpers.SetupHandlers(app)
	app.Listen(":8080")
}
