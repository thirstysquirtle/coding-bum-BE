package auth

import (
	"fmt"
	"os"
	"strings"
	"time"
	"wut/db"
	"wut/ent"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var ecret []byte
var isTest bool = false

func SetupCustomAuth(app *fiber.App, dbCl *db.DbClient) {
	GetSecret()

	app.Post("/auth/custom/login", func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err == nil {
			user, err := dbCl.AuthUser(c.Context(), strings.Join(form.Value["email"], ""), strings.Join(form.Value["password"], ""))
			if err != nil {
				fmt.Println("auth error: %w", err)
				return c.JSON(fiber.Map{"error": err.Error()})
			}
			fmt.Println("gud")
			cookie, err := createCookieJWT(user)
			if err != nil {
				fmt.Println("cookie creation error: %w", err)
				return c.JSON(fiber.Map{"error": err.Error()})
			}
			c.Cookie(cookie)
			return c.JSON("authentication succesful")

		}
		return c.SendString(fmt.Sprintf("Error: %w", err))
	})

}

func createCookieJWT(user *ent.User) (*fiber.Cookie, error) {
	today := time.Now()
	expireDate := today.Add(time.Hour * 24 * 15)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(today),
		ExpiresAt: jwt.NewNumericDate(expireDate),
		Subject:   fmt.Sprint(user.ID)})
	signedToken, err := token.SignedString(ecret)
	if err != nil {
		fmt.Println("Error Signing Token: %w", err)
		return nil, fmt.Errorf("server couldn't sign token, try again later or contact us by email")
	}

	return &fiber.Cookie{
		Name:     "ses",
		Value:    signedToken,
		Expires:  expireDate,
		HTTPOnly: !isTest,
		Secure:   !isTest,
	}, nil
}

func GetSecret() {
	if os.Getenv("IsTest") == "true" {
		os.Setenv("secret1", "testSecret")
		isTest = true
	}

	secretKeyStr, found := os.LookupEnv("secret1")
	if !found {
		panic("No secret Env Variable Found!")
	}
	ecret = []byte(secretKeyStr)
}
