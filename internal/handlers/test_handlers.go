package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/BookBits/bookbits-editor/internal/helpers/renderer"
	"github.com/BookBits/bookbits-editor/templates/components"
	"github.com/BookBits/bookbits-editor/templates/views"
)

var count uint = 0;

func TestPage(c fiber.Ctx) error {
	return renderer.RenderTempl(c, views.TestPage(count))
}

func TestIncrement(c fiber.Ctx) error {
	count += 1
	return renderer.RenderTempl(c, components.Counter(count))
}
