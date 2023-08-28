package main

import (
	"fmt"
	"wut/auth"
	"wut/db"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	// fmt.Println(wrd.test())
	// secretKey := ifDevEnv()
	dbClient := db.SetupDBConnection()
	defer (*dbClient).Client.Close()
	auth.SetupLoginRoutes(app, dbClient)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./index.html")
	})

	app.Listen(":3000")
	fmt.Println("exit")
}
