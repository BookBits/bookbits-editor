package middlewares

import (
	"strings"

	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		state.User = models.User{}
		return c.Next()	
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.Set("Authorization", "")
		c.ClearCookie()
		return c.Redirect().To("/login")
	}

	accessToken := parts[1]

	claims, err := models.ValidateToken(accessToken, state.Vars.JWTSecretKey)

	if err != nil {
		state.User = models.User{}
		return c.Next()
	}

	user, err := models.GetUserByID(claims.UserID, state.DB)
	if err != nil {
		return c.SendStatus(500)
	}

	state.User = user
	return c.Next()
}
