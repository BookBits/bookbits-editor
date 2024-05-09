package handlers

import (
	"github.com/BookBits/bookbits-editor/internal/helpers/renderer"
	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/BookBits/bookbits-editor/templates/views/app"
	"github.com/gofiber/fiber/v3"
)

func GetUsers(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	users, err := models.GetUsers(state.DB)

	if err != nil {
		return c.Status(500).SendString("Error Fetching Users")
	}
	return renderer.RenderTempl(c, app.UserList(users))
}
