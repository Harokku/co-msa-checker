package router

import (
	"co-msa-checker/handler"
	"github.com/gofiber/fiber/v2"
)

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

	// -------------------------
	// Fleet
	// -------------------------
	fleet := v1.Group("/fleet")
	fleet.Get("/", handler.Fleet{}.GetAll)
	fleet.Get("/:id", handler.Fleet{}.GetById)

	// -------------------------
	// Info List
	// -------------------------
	infolist := v1.Group("/info")
	infolist.Get("/all/:id", handler.InfoList{}.GetllByMsaId)
	infolist.Post("/", handler.InfoList{}.PostInfo)
	infolist.Get("/updates/all/:id", handler.InfoList{}.GetAllUpdatesByInfoId)
	infolist.Post("/updates/", handler.InfoList{}.PostUpdate)
}
