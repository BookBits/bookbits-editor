package middlewares

import (
	"errors"
	"strings"
	"time"

	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/gofiber/fiber/v3"
)

func tryRefresh(c fiber.Ctx) (models.User, string, string, error) {
	refreshToken := c.Cookies("refreshToken")

	if refreshToken == "" {
		return models.User{}, "", "", errors.New("no refresh token")
	}

	state := c.Locals("state").(*models.AppState)
	claims, err := models.ValidateToken(refreshToken, state.Vars.JWTSecretKey)

	if err != nil {
		return models.User{}, "", "", err
	}

	user, err := models.GetUserByID(claims.UserID, state.DB)
	if err != nil {
		return models.User{}, "", "", err
	}

	accessToken, refreshToken, err := user.GenerateTokens(state.Vars)
	if err != nil {
		return models.User{}, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func AuthMiddleware(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		c.ClearCookie()
		return c.Redirect().To("/login")
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.Set("Authorization", "")
		c.ClearCookie()
		return c.Redirect().To("/login")
	}

	state := c.Locals("state").(*models.AppState)
	accessToken := parts[1]

	claims, err := models.ValidateToken(accessToken, state.Vars.JWTSecretKey)

	if err != nil {
		if err.Error() == "invalid token" {
			user, accessToken, refreshToken, err := tryRefresh(c)
			if err != nil {
				return c.Redirect().To("/login")
			}

			state.User = user
			c.Cookie(&fiber.Cookie{
				Name: "accessToken",
				Value: accessToken,
				SameSite: "strict",
				HTTPOnly: true,
				Secure: true,
				SessionOnly: true,
				Expires: time.Now().Add(time.Hour * 1),
			})
			c.Cookie(&fiber.Cookie{
				Name: "refreshToken",
				Value: refreshToken,
				SameSite: "strict",
				HTTPOnly: true,
				Secure: true,
				Expires: time.Now().Add(time.Hour * 24 * 7),
			})

			return c.Next()
		} else {
			return c.Redirect().To("/login")
		}
	}

	user, err := models.GetUserByID(claims.UserID, state.DB)
	if err != nil {
		return c.Redirect().To("/login")
	}

	state.User = user
	return c.Next()
}
