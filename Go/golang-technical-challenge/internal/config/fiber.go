package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(v *viper.Viper) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      v.GetString("APP_NAME"),
		ErrorHandler: NewErrorHandler(),
		Prefork:      v.GetBool("WEB_PREFORK"),
	})

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
}
