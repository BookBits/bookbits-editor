package models

import (
	"github.com/go-redis/cache/v9"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type AppState struct {
	DB *gorm.DB
	Cache *cache.Cache
	Vars EnvVars
	User User
	GitClient GitClient
}

func WithAppState(db *gorm.DB, vars EnvVars, gc GitClient, cache *cache.Cache) func (c fiber.Ctx) error {
	state := AppState{DB: db, Vars: vars, User: User{}, GitClient: gc, Cache: cache}
	return func (c fiber.Ctx) error {
		c.Locals("state", &state)
		return c.Next()
	}
}
