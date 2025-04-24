package http

import (
	"github.com/gofiber/fiber/v2"
	"golang-rest/internal/core/ports"
	"golang-rest/internal/core/services"
	"golang-rest/internal/infrastructure/middleware"
)

func Setup(app *fiber.App, userRepository ports.UserRepositoryInterface) {
	userHandlerService := services.NewUserHandlerService(userRepository)

	app.Post("/register", func(ctx *fiber.Ctx) error {
		return userHandlerService.RegisterUser(ctx)
	})

	app.Post("/login", func(ctx *fiber.Ctx) error {
		return userHandlerService.LoginUser(ctx)
	})

	app.Use(middleware.Protected())

	app.Get("/users", middleware.Protected(), func(ctx *fiber.Ctx) error {
		return userHandlerService.GetAllUsers(ctx)
	})

	app.Get("/users/:id", middleware.Protected(), func(ctx *fiber.Ctx) error {
		return userHandlerService.GetUserByID(ctx)
	})

	app.Put("/users/:id", middleware.Protected(), func(ctx *fiber.Ctx) error {
		return userHandlerService.UpdateUserByID(ctx)
	})

	app.Delete("/users/:id", middleware.Protected(), func(ctx *fiber.Ctx) error {
		return userHandlerService.DeleteUserByID(ctx)
	})
}
