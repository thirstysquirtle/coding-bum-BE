package main

import (
	"wut/db"

	"fmt"
	"wut/stripePi"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	var err error
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://localhost:3000",
		AllowCredentials: true,
	}))
	// fmt.Println(wrd.test())
	// secretKey := ifDevEnv()
	db.SetupDBConnection()
	// defer db.PostgresPool.Close()
	// auth.SetupLoginRoutes(app, dbClient)
	stripePi.SetupStripeRoutes(app)
	app.Static("/", "svelte_build/")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("svelte_build/test-auth.html")
	})
	// err := db.Db.CreateUser(context.Background(), sqlc.CreateUserParams{Email: "test", Username: "test", OrderNum: "test"})
	// if err != nil {
	// 	panic(err)
	// }

	err = app.Listen(":3000")
	fmt.Println("wut: %w", err)
	fmt.Println("exit")
	// fmt.Println(time.Now().Add(time.Hour * 24 * 14))
	db.PostgresPool.Close()
}
