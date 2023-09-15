package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
	"wut/db"
	"wut/paymentFlow"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	setDevEnv()
	var err error
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://localhost:3000",
		AllowCredentials: true,
	}))
	db.SetupDBConnection()
	// auth.SetupAuthRoutes(app, dbClient)
	paymentFlow.SetupPaymentFlow(app)

	// app.Static("/", "svelte_build/")

	app.Get("/super-cool-kids-list", requireAuth, func(c *fiber.Ctx) error {
		userPos, err := db.Db.GetUserPos(c.Context(), int32(c.Locals("uid").(int)))
		if err != nil {
			fmt.Println("Error Querying Database: ", err)
			return c.SendStatus(500)
		}
		pageNum := math.Ceil(float64(userPos) / 20)
		users, err := db.Db.GetPageOf20Users(c.Context(), pageNum)
		if err != nil {
			fmt.Println("Error Querying Database: ", err)
			return c.SendStatus(500)
		}
		return c.JSON(fiber.Map{"you": userPos, "coolKids": users})
	})

	err = app.Listen(":3000")
	fmt.Println("wut: %w", err)
	fmt.Println("exit")
	db.PostgresPool.Close()
}

func JwtStringIsValid(jwtString string) (bool, error) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("secret1")), nil
	})

	if token.Valid {
		fmt.Println("You look nice today")
		return true, nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return false, fmt.Errorf("token malformed")
	} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return false, fmt.Errorf("invalid signature")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return false, fmt.Errorf("token time issue")
	} else {
		return false, err
	}
}

func setDevEnv() {
	if os.Getenv("IsTest") == "true" {
		os.Setenv("secret1", "testSecret")
	}

	_, found := os.LookupEnv("secret1")
	if !found {
		panic("No secret Env Variable Found!")
	}
}

func requireAuth(c *fiber.Ctx) error {
	loggedIn := c.Cookies("loggedIn")
	jwtString := c.Cookies("ses")
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("secret1")), nil
	})

	if err != nil || loggedIn == "" {
		// fmt.Println(err, c.Cookies("ses"))
		clearCookies(c, "ses", "loggedIn")

		return c.Redirect("http://localhost:5173/login")
	}
	if token.Valid {
		sub, _ := token.Claims.GetSubject()

		userId, err := strconv.Atoi(sub)
		if err != nil {
			fmt.Println("Token Subject Error: ", token.Raw)
			// c.ClearCookie("ses", "loggedIn")
			clearCookies(c, "ses", "loggedIn")
			return c.SendString(`{"redirecasdt":"/lsogin"}`)
		}
		c.Locals("uid", userId)
		return c.Next()
	}
	fmt.Println("Unexpected results, checkLoginCookies")
	// c.ClearCookie("ses", "loggedIn")
	clearCookies(c, "ses", "loggedIn")
	return c.SendString(`{"rediasdrect":"/logian"}`)
}

func clearCookies(c *fiber.Ctx, key ...string) {
	// expireDate :=
	for _, k := range key {
		c.Cookie(&fiber.Cookie{
			Name:     k,
			Value:    "0",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
			Secure:   true,
		})
	}
}
