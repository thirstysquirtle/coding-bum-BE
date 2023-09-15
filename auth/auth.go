package auth

import (
	// 	"wut/xauth/customAuth"

	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func SetupAuthRoutes(app *fiber.App) {
	SetupCustomAuth(app)
	// app.Get("/self-ranking", func(c *fiber.Ctx) error {
	// 	db.Db.GetUserPos(c.Context())
	// 	return c.JSON(fiber.Map{"position": position})
	// })
	// fmt.Println("Test")
}

func CreateLoginCookies(userId int32) ([]*fiber.Cookie, error) {
	today := time.Now()
	expireDate := today.Add(time.Hour * 24 * 15)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(today),
		ExpiresAt: jwt.NewNumericDate(expireDate),
		Subject:   fmt.Sprint(userId)})
	signedToken, err := token.SignedString([]byte(os.Getenv("secret1")))
	if err != nil {
		fmt.Println("Error Signing Token: ", err)
		return nil, fmt.Errorf("server couldn't sign token, try again later or contact us by email")
	}

	return []*fiber.Cookie{
		{
			Name:     "ses",
			Value:    signedToken,
			Expires:  expireDate,
			HTTPOnly: true,
			Secure:   true,
		},
		{
			Name:     "loggedIn",
			Value:    "true",
			Expires:  expireDate,
			HTTPOnly: false}}, nil
}
