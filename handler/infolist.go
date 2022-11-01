package handler

import (
	"co-msa-checker/database"
	"co-msa-checker/utils"
	"github.com/gofiber/fiber/v2"
)

type InfoList struct {
}

func (f InfoList) GetllByMsaId(ctx *fiber.Ctx) error {
	var (
		res []database.Info
		err error
	)

	// Get data from db
	res, err = database.InfoGetAllByMsaId(ctx.Params("id"))
	if err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved all info based on MSA Id",
		"retrieved": len(res),
		"data":      res,
	})
}

/*func (f InfoList) GetAllWithUpdates(ctx *fiber.Ctx) error {
	var (
		res []database.Info
		err error
	)
}*/

func (f InfoList) GetAllUpdatesByInfoId(ctx *fiber.Ctx) error {
	var (
		res []database.Update
		err error
	)

	//Get data from db
	res, err = database.InfoGetAllUpdatesByInfoId(ctx.Params("id"))
	if err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Retrieved all updates based on info id",
		"retrieved": len(res),
		"data":      res,
	})
}

func (f InfoList) PostUpdate(ctx *fiber.Ctx) error {
	var (
		res database.Update
		err error
	)
	u := new(database.Update)

	// Read data from req body
	if err := ctx.BodyParser(u); err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// Create new Update entry in DB
	res, err = database.NewUpdate(*u)
	if err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Added new update",
		"retrieved": 1,
		"data":      res,
	})
}

func (f InfoList) PostInfo(ctx *fiber.Ctx) error {
	var (
		res database.Info
		err error
	)
	u := new(database.Info)

	// Read data from req body
	if err := ctx.BodyParser(u); err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// Create new Info entry in DB
	res, err = database.NewInfo(*u)
	if err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.JSON(fiber.Map{
		"status":    "success",
		"message":   "Added new info",
		"retrieved": 1,
		"data":      res,
	})
}
