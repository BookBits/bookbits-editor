package middlewares

import (

	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func AuthMiddleware(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)
	accessToken := c.Cookies("accessToken")

	if accessToken == "" {
		state.User = models.User{}
		return c.Next()	
	}

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

func AuthOnlyRoute(c fiber.Ctx) error {
	state := c.Locals("state").(*models.AppState)

	if state.User.ID == uuid.Nil {
		return c.Redirect().To("/")
	}

	return c.Next()
}
