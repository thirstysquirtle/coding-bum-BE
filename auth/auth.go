package auth

import (
	// 	"wut/xauth/customAuth"

	"wut/db"

	"github.com/gofiber/fiber/v2"
)

func SetupLoginRoutes(app *fiber.App, dbCl *db.DbClient) {
	SetupCustomAuth(app, dbCl)
}
