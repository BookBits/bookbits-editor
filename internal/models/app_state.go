package models

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type AppState struct {
	DB *gorm.DB
	Vars EnvVars
	User User
}

func WithAppState(db *gorm.DB, vars EnvVars) func (c fiber.Ctx) error {
	state := AppState{DB: db, Vars: vars, User: User{}}
	return func (c fiber.Ctx) error {
		c.Locals("state", &state)
		return c.Next()
	}
}
