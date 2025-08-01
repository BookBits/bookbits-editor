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
		return c.Status(500).SendString("Error while trying to register user")
	}

	users, err := models.GetUsers(db)

	if err != nil {
		return c.Status(500).SendString("Error while fetching users")
	}
	
	token := csrf.TokenFromContext(c)

	return renderer.RenderTempl(c, app.UserList(users, token))
}

func DeleteUser(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	idVal := c.Params("uid")
	id, err := uuid.Parse(idVal)

	if err != nil {
		return c.Status(400).SendString("Invalid User ID")
	}

	delErr := models.DeleteUserByID(id, state.DB)
	if delErr != nil {
		return c.Status(500).SendString("Error while trying to delete user")
	}

	users, err := models.GetUsers(state.DB)
	if err != nil {
		return c.Status(500).SendString("Error while fetching users")
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
		return c.Status(400).SendString("Invalid user type provided")
	}

	newTypeParsed := models.UserType(newType)
	updateErr := state.DB.Model(&models.User{ID: id}).Update("type", newTypeParsed).Error

	if updateErr != nil {
		return c.Status(500).SendString("Error while trying to update user information")
	}

	user, err := models.GetUserByID(id, state.DB)
	if err != nil {
		return c.Status(500).SendString("Error while fetching user data")
	}
	token := csrf.TokenFromContext(c)
	return renderer.RenderTempl(c, app.UserTile(user, token))
}

func ChangePassword(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	currentPassword := c.FormValue("current-password")
	newPassword := c.FormValue("new-password")

	err := state.User.UpdatePassword(currentPassword, newPassword, state.DB)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	c.ClearCookie()
	c.Set("HX-Redirect", "/login")
	return c.SendStatus(200)
}

func ChangePasswordRandom(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	idVal := c.Params("uid")
	id, err := uuid.Parse(idVal)
	if err != nil {
		return c.Status(400).SendString("Invalid user")
	}

	user, err := models.GetUserByID(id, state.DB)
	if err != nil {
		log.Fatal(err)
		return c.SendStatus(500)
	}

	newPass, err := user.UpdatePasswordRandom(state.DB)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).SendString(newPass)
}
