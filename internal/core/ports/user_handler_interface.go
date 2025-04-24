package ports

import (
	"github.com/gofiber/fiber/v2"
)

type UserHandlerInterface interface {
	RegisterUser(ctx *fiber.Ctx) error
	LoginUser(ctx *fiber.Ctx) error
	GetAllUsers(ctx *fiber.Ctx) error
	GetUserByID(ctx *fiber.Ctx) error
	UpdateUserByID(ctx *fiber.Ctx) error
	DeleteUserByID(ctx *fiber.Ctx) error
}
