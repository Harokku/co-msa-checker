package router

import "github.com/gofiber/fiber/v2"

func Setup(app *fiber.App) {
	// -------------------------
	// Grouping and versioning
	// -------------------------

	api := app.Group("/api")

	v1 := api.Group("/v1", func(ctx *fiber.Ctx) error {
		ctx.Set("Version", "v1")
		return ctx.Next()
	})

	// -------------------------
	// Versions landing
	// -------------------------

	v1.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("API version 1 root")
	})
}
