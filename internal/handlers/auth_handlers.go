package handlers

import (
	"time"

	"github.com/BookBits/bookbits-editor/internal/helpers/renderer"
	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/BookBits/bookbits-editor/templates/views"
	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func IndexPage(c fiber.Ctx) error {
	return renderer.RenderTempl(c, views.IndexPage())
}

func LoginPage(c fiber.Ctx) error {
	return renderer.RenderTempl(c, views.LoginPage())
}

type loginResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt time.Time `json:"expires_at"`
}

func Login(c fiber.Ctx) error {
	userEmail := c.FormValue("user-email")
	password := c.FormValue("user-password")

	state := c.Locals("state").(*models.AppState)

	var user models.User;
	err := state.DB.Where("email = ?", userEmail).First(&user).Error

	if err != nil {
		return c.Status(422).SendString("Invalid Email Address")
	}

	validateErr := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if validateErr != nil {
		return c.SendStatus(fiber.ErrUnauthorized.Code)
	}

	accessToken, refreshToken, err := user.GenerateTokens(state.Vars)
	if err != nil {
		return c.SendStatus(500)
	}

	c.Cookie(&fiber.Cookie{
		Name: "refreshToken",
		Value: refreshToken,
		SameSite: "strict",
		HTTPOnly: true,
		Secure: true,
		Expires: time.Now().Add(time.Hour * 24 * 7),
	})

	loginResponse := loginResponse{
		AccessToken: accessToken,
		ExpiresAt: time.Now().Add(time.Second * 120),
	}

	return c.JSON(loginResponse)
}
