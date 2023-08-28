package auth

import (
	"wut/db"

	"github.com/gofiber/fiber/v2"
)

func SetupCustomAuth(app *fiber.App, dbCl *db.DbClient) {
	app.Post("/auth/custom/register", func(c *fiber.Ctx) error {

		user, _ := dbCl.CreateUser(c.Context(), c.FormValue("email"), c.FormValue("password"))
		return c.JSON(user)
	})

	// app.Get("/jwt", func(c *fiber.Ctx) error {
	// 	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 		"iss":       "codingBum",
	// 		"timeStamp": time.Now(),
	// 	}).SignedString(secretKey)
	// 	if err != nil {
	// 		fmt.Println("Jwt signing Error:", err)
	// 		return c.Status(500).SendString("Server Error")
	// 	}
	// 	return c.SendString(jwtToken)
	// })

}
