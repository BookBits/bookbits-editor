package handlers

import (
	"github.com/BookBits/bookbits-editor/internal/helpers/renderer"
	"github.com/BookBits/bookbits-editor/templates/views"
	"github.com/gofiber/fiber/v3"
)

func IndexPage(c fiber.Ctx) error {
	return renderer.RenderTempl(c, views.IndexPage())
}
