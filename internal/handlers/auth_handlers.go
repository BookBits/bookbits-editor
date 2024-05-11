package handlers

import (
	"log"
	"time"

	"github.com/BookBits/bookbits-editor/internal/helpers/renderer"
	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/BookBits/bookbits-editor/templates/views"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/google/uuid"
)

type csrfRes struct {
	CsrfToken string `json:"csrfToken"`
}

func GetCSRF(c fiber.Ctx) error {
	token := csrf.TokenFromContext(c)
	res := csrfRes{CsrfToken: token}
	return c.JSON(res)
}

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
		log.Fatal(err)
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
		log.Fatal(err)
		return c.SendStatus(500)
	}

	c.Cookie(&fiber.Cookie{
		Name: "refreshToken",
		Value: refreshToken,
		SameSite: "Strict",
		Expires: time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
	})
	
	c.Cookie(&fiber.Cookie{
		Name: "accessToken",
		Value: accessToken,
		SameSite: "Strict",
		HTTPOnly: true,
		Expires: time.Now().Add(time.Hour),
	})
	
	loginResponse := loginResponse{
		ExpiresAt: time.Now().Add(time.Hour),
	}
	return c.JSON(loginResponse)
}

func IndexHandler(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	user := state.User

	if user.ID == uuid.Nil {
		token := csrf.TokenFromContext(c)
		return renderer.RenderTempl(c, views.IndexPage(token))
	}

	hxReq := c.Get("HX-Request")
	
	if hxReq == "true" {
		c.Set("HX-Redirect", "/app")
		return c.SendStatus(200)
	}

	return c.Redirect().To("/app")
}

func LoginPageHandler(c fiber.Ctx) error {
	token := csrf.TokenFromContext(c)
	return renderer.RenderTempl(c, views.LoginPage(token))
}

type loginResponse struct {
	ExpiresAt time.Time `json:"expires_at"`
}

func Login(c fiber.Ctx) error {
	userEmail := c.FormValue("user-email")
	password := c.FormValue("user-password")

	state := c.Locals("state").(*models.AppState)

	var user models.User;
	err := state.DB.Where("email = ?", userEmail).First(&user).Error

	if err != nil {
		log.Fatal(err)
		return c.Status(422).SendString("Invalid Email Address")
	}

	validateErr := user.Validate(password)
	if validateErr != nil {
		log.Fatal(validateErr)
		return c.Status(401).SendString("Invalid Password")
	}

	accessToken, refreshToken, err := user.GenerateTokens(state.Vars)
	if err != nil {
		log.Fatal(err)
		return c.Status(500).SendString("Server Error")
	}

	c.Cookie(&fiber.Cookie{
		Name: "refreshToken",
		Value: refreshToken,
		SameSite: "Strict",
		HTTPOnly: true,
		Expires: time.Now().Add(time.Hour * 24 * 7),
	})

	c.Cookie(&fiber.Cookie{
		Name: "accessToken",
		Value: accessToken,
		SameSite: "Strict",
		HTTPOnly: true,
		Expires: time.Now().Add(time.Hour),
	})
	
	loginResponse := loginResponse{
		ExpiresAt: time.Now().Add(time.Hour),
	}
	return c.JSON(loginResponse)
}

func Logout(c fiber.Ctx) error {
	c.ClearCookie("accessToken")
	c.ClearCookie("refreshToken")

	c.Set("HX-Redirect", "/")
	return c.SendStatus(200)
}
