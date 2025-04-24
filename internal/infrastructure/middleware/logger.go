package middleware

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func Logger() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()
		err := ctx.Next()
		log.Printf("%s %s %v\n", ctx.Method(), ctx.Path(), time.Since(start))
		return err
	}
}
