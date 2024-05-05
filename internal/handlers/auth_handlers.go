package handlers

import (
	"github.com/BookBits/bookbits-editor/internal/helpers/renderer"
	"github.com/BookBits/bookbits-editor/templates/views"
	"github.com/gofiber/fiber/v3"
)

func IndexPage(c fiber.Ctx) error {
	return renderer.RenderTempl(c, views.IndexPage())
}

func LoginPage(c fiber.Ctx) error {
	return renderer.RenderTempl(c, views.LoginPage())
}

func Login(c fiber.Ctx) error {
	c.Set("HX-Redirect", "/dashboard")
	c.Cookie(&fiber.Cookie{
		Name: "accessToken",
		Value: "some token",
		Secure: true,
		SameSite: "strict",
		HTTPOnly: true,
		SessionOnly: true,
	})
	return c.SendString("")
}
