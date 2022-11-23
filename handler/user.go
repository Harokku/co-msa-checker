package handler

import (
	"co-msa-checker/database"
	"co-msa-checker/utils"
	"github.com/gofiber/fiber/v2"
)

type User struct {
}

// Login get password prom post, check if user exist in db and return userdata in a JSON
func (u User) Login(ctx *fiber.Ctx) error {
	var (
		res database.User
		err error
	)
	bodyUser := new(database.User)

	// Read data from req body
	if err = ctx.BodyParser(bodyUser); err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// Look for password/user match in db
	res, err = database.GetUserDataByPassword(bodyUser.Password)
	if err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "User authorized",
		"data":    res,
	})
}
