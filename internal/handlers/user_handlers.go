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

func GetUsers(c fiber.Ctx) error {
	token := csrf.TokenFromContext(c)
	state := c.Locals("state").(*models.AppState)
	users, err := models.GetUsers(state.DB)

	if err != nil {
		return c.Status(500).SendString("Error Fetching Users")
	}
	return renderer.RenderTempl(c, app.UserList(users, token))
}

func RegisterUser(c fiber.Ctx) error {
	username := c.FormValue("username")
	user_email := c.FormValue("user-email")
	password := c.FormValue("user-password")
	user_type := c.FormValue("user-type")

	state := c.Locals("state").(*models.AppState)
	db := state.DB

	err := models.CreateUserWithPassword(username, user_email, password, models.UserType(user_type), db)

	if err != nil {
		log.Fatal(err)
		return c.Status(500).SendString("Unable to create user")
	}

	users, err := models.GetUsers(db)

	if err != nil {
		return c.Status(500).SendString("Unable to fetch users")
	}
	
	token := csrf.TokenFromContext(c)

	return renderer.RenderTempl(c, app.UserList(users, token))
}

func DeleteUser(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	idVal := c.Params("uid")
	id, err := uuid.Parse(idVal)

	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString("Invalid User ID")
	}

	delErr := models.DeleteUserByID(id, state.DB)
	if delErr != nil {
		return c.SendStatus(500)
	}

	users, err := models.GetUsers(state.DB)
	if err != nil {
		return c.Status(500).SendString("Unable to fetch users")
	}
	token := csrf.TokenFromContext(c)

	return renderer.RenderTempl(c, app.UserList(users, token))
}

func UpdateUserType(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	idVal := c.Params("uid")
	id, err := uuid.Parse(idVal)

	newType := c.FormValue("new-type")
	if newType == "" {
		return c.Status(400).SendString("Invalid user type")
	}

	newTypeParsed := models.UserType(newType)
	updateErr := state.DB.Model(&models.User{ID: id}).Update("type", newTypeParsed).Error

	if updateErr != nil {
		return c.SendStatus(500)
	}

	user, err := models.GetUserByID(id, state.DB)
	if err != nil {
		return c.Status(500).SendString("Error Fetching User")
	}
	token := csrf.TokenFromContext(c)
	return renderer.RenderTempl(c, app.UserTile(user, token))
}
