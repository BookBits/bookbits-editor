package handlers

import (
	"github.com/BookBits/bookbits-editor/internal/helpers/renderer"
	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/BookBits/bookbits-editor/templates/views/app"
	"github.com/gofiber/fiber/v3"
)

func AppHomeHandler(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)

	return renderer.RenderTempl(c, app.AppHomePage(state.User))
}
