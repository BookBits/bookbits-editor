package handlers

import (
	"time"

	"github.com/BookBits/bookbits-editor/internal/helpers/renderer"
	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/BookBits/bookbits-editor/templates/views"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RefreshSession(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	refreshToken := c.Cookies("refreshToken")
	hxRequest := c.Get("HX-Request")

	if refreshToken == "" {
		if (hxRequest == "true") {
			c.Set("HX-Redirect", "/login")
			return c.SendStatus(302)
		} else {
			return c.Redirect().To("/login")
		}
	}

	claims, err := models.ValidateToken(refreshToken, state.Vars.JWTSecretKey)
	if err != nil {
		if (hxRequest == "true") {
			c.Set("HX-Redirect", "/login")
			return c.SendStatus(302)
		} else {
			return c.Redirect().To("/login")
		}
	}

	user, err := models.GetUserByID(claims.UserID, state.DB)
	accessToken, refreshToken, err := user.GenerateTokens(state.Vars)

	if err != nil {
		return c.SendStatus(500)
	}

	c.Cookie(&fiber.Cookie{
		Name: "refreshToken",
		Value: refreshToken,
		Secure: true,
		SameSite: "Strict",
		Expires: time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
	})

	response := loginResponse{
		AccessToken: accessToken,
		ExpiresAt: time.Now().Add(time.Minute * 2),
	}
	return c.Status(200).JSON(&response)
}

func IndexHandler(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	user := state.User

	if user.ID == uuid.Nil {
		return renderer.RenderTempl(c, views.IndexPage())
	}

	c.Set("HX-Redirect", "/app")
	return c.SendStatus(200)
}

func LoginPageHandler(c fiber.Ctx) error {
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
		return c.Status(401).SendString("Invalid Password")
	}

	accessToken, refreshToken, err := user.GenerateTokens(state.Vars)
	if err != nil {
		return c.Status(500).SendString("Server Error")
	}

	c.Cookie(&fiber.Cookie{
		Name: "refreshToken",
		Value: refreshToken,
		SameSite: "Strict",
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
