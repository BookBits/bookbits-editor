package handlers

import (
	"github.com/BookBits/bookbits-editor/internal/helpers/renderer"
	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/BookBits/bookbits-editor/templates/views/app"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/google/uuid"
)

func GetProjects(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	projects, err := state.User.GetProjects(state.DB)

	if err != nil {
		return c.Status(500).SendString("Unable to fetch projects")
	}

	csrfToken := csrf.TokenFromContext(c)

	return renderer.RenderTempl(c, app.ProjectsList(projects, state.User, csrfToken))
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

func DeleteProject(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	if state.User.Type == models.UserTypeWriter {
		return c.SendStatus(401)
	}
	projectID, err := uuid.Parse(c.Params("pid"))

	if err != nil {
		return c.Status(500).SendString("Invalid Project ID")
	}

	var project models.Project
	if err := state.DB.First(&project, projectID).Error;err != nil {
		return c.Status(500).SendString("Invalid Project")
	}

	if err := project.Delete(state); err != nil {
		log.Error(err)
		return c.Status(500).SendString("Unable to Delete Project")
	}

	return GetProjects(c)
}
