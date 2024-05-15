package helpers

import (
	"github.com/BookBits/bookbits-editor/internal/handlers"
	"github.com/BookBits/bookbits-editor/internal/middlewares"
	"github.com/gofiber/fiber/v3"
)

func SetupHandlers(app *fiber.App) {
	app.Static("/static", "./public/")
	
	app.Get("/test", handlers.TestPage)
	app.Post("/test/increment", handlers.TestIncrement)
	app.Get("/check_db", handlers.CheckDB)

	app.Get("/", handlers.IndexHandler, middlewares.AuthMiddleware)
	app.Post("/refresh", handlers.RefreshSession)
	app.Get("/csrf", handlers.GetCSRF)
	app.Get("/login", handlers.LoginPageHandler)
	app.Post("/login", handlers.Login)
	app.Post("/logout", handlers.Logout)

	app.Get("/app", handlers.AppHomeHandler, middlewares.AuthMiddleware, middlewares.AuthOnlyRoute)
	app.Get("/users", handlers.GetUsers, middlewares.AuthMiddleware, middlewares.AdminOnlyRoute)
	app.Post("/users", handlers.RegisterUser, middlewares.AuthMiddleware, middlewares.AdminOnlyRoute)
	app.Delete("/users/:uid", handlers.DeleteUser, middlewares.AuthMiddleware, middlewares.AdminOnlyRoute)
	app.Patch("/users/:uid/type", handlers.UpdateUserType, middlewares.AuthMiddleware, middlewares.AdminOnlyRoute)
	app.Patch("/users/password", handlers.ChangePassword, middlewares.AuthMiddleware, middlewares.AuthOnlyRoute)
	app.Patch("/users/:uid/password", handlers.ChangePasswordRandom, middlewares.AuthMiddleware, middlewares.AdminOnlyRoute)

	app.Get("/projects", handlers.GetProjects, middlewares.AuthMiddleware, middlewares.AuthOnlyRoute)
	app.Post("/projects", handlers.CreateProject, middlewares.AuthMiddleware, middlewares.AuthOnlyRoute)
}
