package handler

import (
	"co-msa-checker/database"
	"co-msa-checker/utils"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
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

func (u User) CreateUsersFromXls(ctx *fiber.Ctx) error {
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

	// Read uploaded file content and build users table
	newUsers, err := database.Xls{}.BuildUsers(readedFile)
	if err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	err = database.BulkCreate(newUsers)
	if err != nil {
		utils.Err(err)
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.SendFile("Utenti.xlsx")
	//return ctx.SendStatus(fiber.StatusOK)

}
