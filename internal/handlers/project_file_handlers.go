package handlers

import (
	"fmt"
	"strconv"
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

	files, err := project.GetFiles(state.DB, state.User);
	if err != nil {
		return c.Status(500).SendString("Error while trying to fetch files for project")
	}
	
	csrfToken := csrf.TokenFromContext(c)
	content := app.ProjectFilesSection(csrfToken, files, project, state.User)
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

	files, err := project.GetFiles(state.DB, state.User)
	if err != nil {
		return c.Status(500).SendString("Error while trying to fetch files for the project. Please Refresh")
	}
	token := csrf.TokenFromContext(c)

	return renderer.RenderTempl(c, app.ProjectFilesList(files, token, state.User))
}

func AddReviewer(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	fileID, err := uuid.Parse(c.Params("fid"))

	if err != nil {
		log.Error(err)
		return c.Status(400).SendString("Trying to modify invalid file")
	}

	var file models.ProjectFile
	if err := state.DB.Preload("Reviewers").First(&file, fileID).Error; err != nil {
		log.Error(err)
		return c.Status(400).SendString("Trying to modify invalid file")
	}
	
	reviewerEmail := c.FormValue("add-reviewer-email")
	if reviewerEmail == "" {
		return c.Status(400).SendString("No Email Address provided for reviewer")
	}
	
	reviewer, err := models.GetUserByEmail(reviewerEmail, state.DB)
	if err != nil {
		return c.Status(400).SendString("No User found for the Email Address Provided")
	}

	if file.EditorID == reviewer.ID {
		return c.Status(400).SendString("User already assigned as editor")
	}

	if err := state.DB.Model(&file).Association("Reviewers").Append(&reviewer); err != nil {
		return c.Status(500).SendString("Couldn't assign reviewer. Please try again.")
	}
	token := csrf.TokenFromContext(c)
	return renderer.RenderTempl(c, app.ReviewersList(token, file))
}

func RemoveReviewer(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	fileID, err := uuid.Parse(c.Params("fid"))

	if err != nil {
		log.Error(err)
		return c.Status(400).SendString("Trying to modify invalid file")
	}

	var file models.ProjectFile
	if err := state.DB.Preload("Reviewers").First(&file, fileID).Error; err != nil {
		log.Error(err)
		return c.Status(400).SendString("Trying to modify invalid file")
	}
	
	reviewerID, err := uuid.Parse(c.Params("reviewerId"))

	if err != nil {
		log.Error(err)
		return c.Status(400).SendString("Trying to remove invalid reviewer")
	}

	if err := state.DB.Model(&file).Association("Reviewers").Delete(&models.User{ID: reviewerID});err != nil {
		return c.Status(500).SendString("Error while trying to remove reviewer")
	}
	token := csrf.TokenFromContext(c)
	return renderer.RenderTempl(c, app.ReviewersList(token, file))
}

func AssignEditor(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	fileID, err := uuid.Parse(c.Params("fid"))

	if err != nil {
		return c.Status(400).SendString("Trying to modify invalid file")
	}

	var file models.ProjectFile
	if err := state.DB.Preload("Reviewers").Preload("Editor").Preload("Creator").First(&file, fileID).Error; err != nil {
		return c.Status(400).SendString("Trying to modify invalid file")
	}
	
	editorEmail := c.FormValue("editor-email")
	if editorEmail == "" {
		return c.Status(400).SendString("No Email Address provided for editor")
	}
	
	editor, err := models.GetUserByEmail(editorEmail, state.DB)
	if err != nil {
		return c.Status(400).SendString("No User found for the Email Address Provided")
	}
	
	if err := state.DB.Model(&file).Association("Editor").Append(&editor); err != nil {
		return c.Status(500).SendString("Error while trying to assign Editor. Please try again.")
	}
	csrfToken := csrf.TokenFromContext(c)
	return renderer.RenderTempl(c, app.ProjectFileListElement(file, csrfToken, state.User))
}

func DeleteFile(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	
	if state.User.Type == models.UserTypeWriter {
		return c.SendStatus(401)
	}

	fileID, err := uuid.Parse(c.Params("fid"))

	if err != nil {
		return c.Status(400).SendString("Trying to modify invalid file")
	}
	
	var file models.ProjectFile
	if err := state.DB.Preload("Editor").Preload("Reviewers").Preload("Creator").Preload("Project").First(&file, fileID).Error; err != nil {
		return c.Status(400).SendString("Trying to modify invalid file")
	}

	if err := file.Delete(state); err != nil {
		log.Error(err)
		return c.Status(500).SendString("Error while trying to Delete File")
	}

	updatedFiles, err := file.Project.GetFiles(state.DB, state.User)
	if err != nil {
		return c.Status(500).SendString("Unable to fetch updated files.")
	}

	csrfToken := csrf.TokenFromContext(c)

	return renderer.RenderTempl(c, app.ProjectFilesList(updatedFiles, csrfToken, state.User))
}

func EditFile(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	fileID, err := uuid.Parse(c.Params("fid"))

	if err != nil {
		return c.Status(400).SendString("Trying to modify invalid file")
	}
	
	var file models.ProjectFile
	if err := state.DB.Preload("Editor").Preload("Reviewers").Preload("Creator").Preload("Project").First(&file, fileID).Error; err != nil {
		return c.Status(400).SendString("Trying to modify invalid file")
	}

	lockUserID, err := file.IsLocked(state)
	if err != nil {
		log.Error(err)
		return c.Status(500).SendString("Error while trying to open file. Please try again.")
	}

	if !(lockUserID == uuid.Nil || lockUserID == state.User.ID) {
		return c.Status(400).SendString("File is locked for editing. Please try later or use the view option.")
	}

	lock, err := file.LockFile(state)
	
	if err != nil {
		log.Error(err)
		return c.Status(500).SendString("Error while trying to open file. Please try again.")
	}
	
	fileContents, err := file.GetContents(state)
	if err != nil {
		return c.Status(500).SendString("Unable to fetch contents for the file")
	}

	c.Cookie(&fiber.Cookie{
		Name: "File-Lock-Expire",
		Value: lock.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
		SessionOnly: true,
		SameSite: "Strict",
	})

	csrfToken := csrf.TokenFromContext(c)
	content := app.Editor(file, fileContents, csrfToken)
	return renderer.RenderTempl(c, app.AppHomePage(state.User, csrfToken, fmt.Sprintf("%v | BookBits Editor", file.Name), content))
}

func RefreshLock(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	fileID, err := uuid.Parse(c.Params("fid"))

	if err != nil {
		return c.Status(400).SendString("Trying to modify invalid file")
	}
	
	var file models.ProjectFile
	if err := state.DB.Preload("Editor").Preload("Reviewers").Preload("Creator").Preload("Project").First(&file, fileID).Error; err != nil {
		return c.Status(400).SendString("Trying to modify invalid file")
	}
	
	lockUserID, err := file.IsLocked(state)
	if err != nil {
		log.Error(err)
		return c.Status(500).SendString("Unable to obtain lock data. Please refresh page.")
	}

	if !(lockUserID == uuid.Nil || lockUserID == state.User.ID) {
		return c.SendStatus(401)
	}

	lock, err := file.LockFile(state)
	if err != nil {
		return c.Status(500).SendString("Unable to lock file. Please refresh page")
	}
	
	c.Cookie(&fiber.Cookie{
		Name: "File-Lock-Expire",
		Value: lock.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
		SessionOnly: true,
		SameSite: "Strict",
	})

	return c.SendStatus(200)
}

func SaveFile(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	fileID, err := uuid.Parse(c.Params("fid"))

	if err != nil {
		return c.Status(400).SendString("Trying to modify invalid file")
	}
	
	fileVersion, err := strconv.Atoi(c.Get(fmt.Sprintf("X-File-%v-Version", fileID)))
	
	if err != nil {
		return c.Status(400).SendString("File Version not included with save request")
	}
	
	var file models.ProjectFile
	if err := state.DB.Preload("Editor").Preload("Reviewers").Preload("Creator").Preload("Project").First(&file, fileID).Error; err != nil {
		return c.Status(400).SendString("Trying to modify invalid file")
	}
	
	lockUserID, err := file.IsLocked(state)
	if err != nil {
		log.Error(err)
		return c.Status(500).SendString("Unable to obtain lock data. Please refresh page.")
	}

	if !(lockUserID == uuid.Nil || lockUserID == state.User.ID) {
		return c.Status(401).SendString("File lock for your session is expired. Please refresh page and try again.")
	}

	if uint(fileVersion) != file.Version {
		return c.Status(400).SendString("File Version mismatch. Can't save as it might cause merge conflicts")
	}

	contents := c.FormValue("content")
	newVersion, saveErr := file.Save(state, []byte(contents))

	if saveErr != nil {
		return c.Status(500).SendString("Couldn't Save the file. Try again")
	}

	csrfToken := csrf.TokenFromContext(c)

	newButton := app.EditorSaveAndContinueButton(fileID, newVersion, csrfToken)

	return renderer.RenderTempl(c, newButton)
}

func UnlockFile(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	fileID, err := uuid.Parse(c.Params("fid"))

	if err != nil {
		return c.Status(400).SendString("Trying to modify invalid file")
	}
	
	var file models.ProjectFile
	if err := state.DB.Preload("Editor").Preload("Reviewers").Preload("Creator").Preload("Project").First(&file, fileID).Error; err != nil {
		return c.Status(400).SendString("Trying to modify invalid file")
	}

	if err := file.UnlockFile(state); err != nil {
		return c.SendStatus(500)
	}
	return c.SendStatus(200)
}
