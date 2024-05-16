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
	projects, err := state.User.GetProjects(state.DB)

	if err != nil {
		return c.Status(500).SendString("Unable to fetch projects")
	}


	content := app.ProjectsSection(csrfToken, state.User, projects)
	return renderer.RenderTempl(c, app.AppHomePage(state.User, csrfToken, "Dashboard | Bookbits Editor", content))
}
