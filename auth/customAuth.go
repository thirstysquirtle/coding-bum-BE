package auth

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

var isTest bool = false

func SetupCustomAuth(app *fiber.App) {
	if os.Getenv("IsTest") == "true" {
		isTest = true
	}

	// app.Post("/auth/custom/login", func(c *fiber.Ctx) error {
	// 	form, err := c.MultipartForm()
	// 	if err == nil {
	// 		user, err := dbCl.AuthUser(c.Context(), strings.Join(form.Value["email"], ""), strings.Join(form.Value["password"], ""))
	// 		if err != nil {
	// 			fmt.Println("auth error: %w", err)
	// 			return c.JSON(fiber.Map{"error": err.Error()})
	// 		}
	// 		fmt.Println("gud")
	// 		cookie, err := createCookieJWT(user)
	// 		if err != nil {
	// 			fmt.Println("cookie creation error: %w", err)
	// 			return c.JSON(fiber.Map{"error": err.Error()})
	// 		}
	// 		c.Cookie(cookie)
	// 		return c.JSON("authentication succesful")

	// 	}
	// 	return c.SendString(fmt.Sprintf("Error: %w", err))
	// })

}
