package main

import (
	"fmt"
	"wut/auth"
	"wut/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://localhost:3000",
		AllowCredentials: true,
	}))
	// fmt.Println(wrd.test())
	// secretKey := ifDevEnv()
	dbClient := db.SetupDBConnection()
	defer (*dbClient).Client.Close()
	auth.SetupLoginRoutes(app, dbClient)

	app.Static("/", "svelte_build/")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("svelte_build/test-auth.html")
	})

	err := app.Listen(":3000")
	fmt.Println("wut: %w", err)
	fmt.Println("exit")
	// fmt.Println(time.Now().Add(time.Hour * 24 * 14))
}
