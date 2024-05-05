package main

import (
	"fmt"
	"log"

	"github.com/BookBits/bookbits-editor/internal/helpers"
	"github.com/BookBits/bookbits-editor/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
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
	app.Use(models.WithAppState(db))

	//Setup Handlers
	helpers.SetupHandlers(app)

	//Start Server
	log.Fatal(app.Listen(fmt.Sprintf(":%s", vars.Port)))
}
