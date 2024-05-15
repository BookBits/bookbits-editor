package handlers

import (
	"github.com/BookBits/bookbits-editor/internal/helpers/renderer"
	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/BookBits/bookbits-editor/templates/views/app"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func GetProjects(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	projects, err := state.User.GetProjects(state.DB)

	if err != nil {
		return c.Status(500).SendString("Unable to fetch projects")
	}

	return renderer.RenderTempl(c, app.ProjectsList(projects))
}

func CreateProject(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	newProjectName := c.FormValue("new-project-name")

	if newProjectName == "" {
		return c.Status(400).SendString("Project Name cannot be empty")
	}

	if err := models.NewProject(newProjectName, state.DB, state.GitClient, state.User.ID); err != nil {
		log.Error(err)
		return c.Status(500).SendString("Unable to Create New Project. Please try Again")
	}

	return GetProjects(c)
}
