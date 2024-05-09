package handlers

import (
	"github.com/BookBits/bookbits-editor/internal/helpers/renderer"
	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/BookBits/bookbits-editor/templates/views/app"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

func AppHomeHandler(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	csrfToken := csrf.TokenFromContext(c)

	return renderer.RenderTempl(c, app.AppHomePage(state.User, csrfToken))
}
