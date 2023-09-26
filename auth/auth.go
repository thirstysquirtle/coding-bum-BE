package auth

import (
	// 	"wut/xauth/customAuth"

	"encoding/base64"
	"fmt"
	"os"
	"time"
	"wut/db"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SetupAuthRoutes(app *fiber.App) {
	app.Post("/login", func(c *fiber.Ctx) error {
		var loginCreds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&loginCreds); err != nil {
			fmt.Println(loginCreds)

			fmt.Println("1:", err)
			return c.SendStatus(500)
		}

		passAndId, err := db.Db.GetUserPassAndId(c.Context(), loginCreds.Email)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}
		// fmt.Println(passAndId.PassHash)
		// fmt.Println(bcrypt.GenerateFromPassword([]byte(loginCreds.Password), bcrypt.DefaultCost))
		err = bcrypt.CompareHashAndPassword([]byte(passAndId.PassHash), []byte(loginCreds.Password))
		if err != nil {
			fmt.Println(err)
			return c.SendString(`{"error":"error"}`)
		}

		loginCookies, err := CreateLoginCookies(passAndId.ID)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(500)
		}
		for _, cookie := range loginCookies {
			c.Cookie(cookie)
		}
		return c.SendString(`{"go-to":"/"}`)
	})

	app.Post("/forgot-password", func(c *fiber.Ctx) error {

		return c.SendStatus(200)
	})

}

func checkPasswordResetToken(tokenString string) bool {
	decodedString := make([]byte, base64.URLEncoding.DecodedLen(len(tokenString)))
	base64.URLEncoding.Decode(decodedString, []byte(tokenString))

	token, err := jwt.Parse(string(decodedString), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("secret1")), nil
	})
	if token.Valid {
		return true
	} else if err != nil {
		fmt.Println("Bad Password Reset Token")
		return false
	}
	fmt.Println("Unexpected")
	return false
}

func createResetPasswordToken(userId int32) []byte {
	expireDate := time.Now().Add(time.Hour * 4)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expireDate),
		Subject:   fmt.Sprint(userId),
		Issuer:    "rp",
	})
	signedToken, err := token.SignedString([]byte(os.Getenv("secret1")))
	if err != nil {
		fmt.Println("Error Signing Token for password reset: ", err)
	}
	urlToken := make([]byte, base64.URLEncoding.EncodedLen(len(signedToken)))
	base64.URLEncoding.Encode(urlToken, []byte(signedToken))
	return urlToken
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
