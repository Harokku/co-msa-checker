package handler

import (
	"co-msa-checker/database"
	"co-msa-checker/utils"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
)

type Fleet struct {
}

// GetAll retrieve all fleet from db
func (f Fleet) GetAll(ctx *fiber.Ctx) error {
	var (
		res []database.Fleet
		err error
	)

	// Get all fleet from db
	res, err = database.FleetGetAll()
	if err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved all fleet",
		"retrieved": len(res),
		"data":      res,
	})
}

// GetById retrieve selected fleet from db
func (f Fleet) GetById(ctx *fiber.Ctx) error {
	var (
		res database.Fleet
		err error
	)

	// Get fleet from db by id
	res, err = database.FleetGetById(ctx.Params("id"))
	switch err {
	case sql.ErrNoRows:
		return ctx.SendStatus(fiber.StatusNotFound)
	case nil:
		return ctx.JSON(fiber.Map{
			"status":    "success",
			"message":   "Retrieved requested fleet",
			"retrieved": 1,
			"data":      res,
		})
	default:
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
}

func (f Fleet) CreateFleetFromXls(ctx *fiber.Ctx) error {
	var (
		file       *multipart.FileHeader
		readedFile multipart.File
		//writer     io.Writer
		err error
	)

	// Get first file from form field "document":
	file, err = ctx.FormFile("config")
	if err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	readedFile, err = file.Open()

	// Read uploaded file content and build fleet table
	newFleet, err := database.Xls{}.BuildFleet(readedFile)
	if err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	err = database.BulkCreateFleet(newFleet)
	if err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.SendStatus(fiber.StatusOK)
}
