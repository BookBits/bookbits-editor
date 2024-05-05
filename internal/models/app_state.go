package models

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type AppState struct {
	DB *gorm.DB
}

func WithAppState(db *gorm.DB) func (c fiber.Ctx) error {
	state := AppState{DB: db}
	return func (c fiber.Ctx) error {
		c.Locals("state", state)
		return c.Next()
	}
}
