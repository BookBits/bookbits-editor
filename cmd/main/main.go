package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/BookBits/bookbits-editor/internal/helpers"
	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/utils/v2"
)

func main() {
	//Setup Logging
	logW, closeFunc, err := helpers.SetupLogger()

	if err != nil {
		return ;
	}
	defer closeFunc()

	// Create the fiber app
	app := fiber.New()

	// Connect Logger
	app.Use(logger.New(logger.Config{Output: logW}))

	//Fetch Environment Variables
	vars, err := helpers.SetupEnvVars()
	if err != nil {
		log.Fatal(err)
		return
	}

	//setup DB
	db, dbErr := helpers.SetupDB(vars)
	if dbErr != nil {
		log.Fatal(dbErr)
		return
	}

	//Setup AppState
	app.Use(models.WithAppState(db, vars))
	
	app.Use(csrf.New(csrf.Config{
		KeyLookup: "header:X-CSRF-Token",
		CookieName: "csrf_",
		CookieSameSite: "Strict",
		CookieHTTPOnly: true,
		CookieSessionOnly: true,
		Expiration: time.Hour * 1,
		KeyGenerator: utils.UUIDv4,
		Extractor: func(c fiber.Ctx) (string, error) {
			token := c.Get("X-CSRF-Token")
			cookie := c.Cookies("csrf_")
			log.Info(cookie)
			log.Info(token)
			if token == "" {
				return token, errors.New("No token")
			}
			return token, nil
		},
	}))

	//Setup Handlers
	helpers.SetupHandlers(app)

	//Start Server
	log.Fatal(app.Listen(fmt.Sprintf(":%s", vars.Port)))
}
