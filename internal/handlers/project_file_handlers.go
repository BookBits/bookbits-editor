package handlers

import (
	"fmt"
	"strings"

	"github.com/BookBits/bookbits-editor/internal/helpers/renderer"
	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/BookBits/bookbits-editor/templates/views/app"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/google/uuid"
)

func GetFiles(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	projectID, err := uuid.Parse(c.Params("pid"))

	if err != nil {
		return c.Status(500).SendString("Invalid Project ID")
	}

	var project models.Project
	if err := state.DB.First(&project, projectID).Error;err != nil {
		return c.Status(500).SendString("Invalid Project")
	}

	files, err := project.GetFiles(state.DB);
	if err != nil {
		return c.Status(500).SendString("Error while trying to fetch files for project")
	}
	
	csrfToken := csrf.TokenFromContext(c)
	content := app.ProjectFilesSection(csrfToken, files, project)
	return renderer.RenderTempl(c, app.AppHomePage(state.User, csrfToken, fmt.Sprintf("%v | Bookbits Editor", project.Name), content))
}

func NewFile(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	projectID, err := uuid.Parse(c.Params("pid"))

	if err != nil {
		return c.Status(500).SendString("Invalid Project ID")
	}

	var project models.Project
	if err := state.DB.First(&project, projectID).Error;err != nil {
		return c.Status(500).SendString("Invalid Project")
	}

	filename := c.FormValue("new-file-name")
	if filename == "" {
		return c.Status(400).SendString("File Name cannot be empty")
	}

	filename = strings.ReplaceAll(filename, " ", "-")

	if err := project.NewFile(filename, state); err != nil {
		log.Error(err)
		return c.Status(500).SendString("Unable to Create File. Please try Again")
	}

	files, err := project.GetFiles(state.DB)
	if err != nil {
		return c.Status(500).SendString("Error while trying to fetch files for the project. Please Refresh")
	}

	return renderer.RenderTempl(c, app.ProjectFilesList(files))
}
