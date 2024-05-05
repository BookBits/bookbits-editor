package handlers

import (
	"github.com/BookBits/bookbits-editor/internal/helpers/renderer"
	"github.com/BookBits/bookbits-editor/internal/models"
	ErrorResponses "github.com/BookBits/bookbits-editor/internal/models/error_responses"
	"github.com/BookBits/bookbits-editor/templates/components"
	"github.com/BookBits/bookbits-editor/templates/views"
	"github.com/gofiber/fiber/v3"
)

var count uint = 0;

func TestPage(c fiber.Ctx) error {
	return renderer.RenderTempl(c, views.TestPage(count))
}

func TestIncrement(c fiber.Ctx) error {
	count += 1
	return renderer.RenderTempl(c, components.Counter(count))
}

func CheckDB(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	db := state.DB

	var health int;
	err := db.Raw("SELECT 1").Find(&health).Error

	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).SendString(ErrorResponses.DBErrorMessage)
	}

	return c.SendString("OK")
}
