package paymentFlow

import (
	"errors"
	"fmt"
	"strconv"
	"wut/auth"
	"wut/db"
	"wut/sqlc"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/paymentintent"
	"golang.org/x/crypto/bcrypt"
)

func SetupPaymentFlow(app *fiber.App) {
	// This is your test secret API key.
	stripe.Key = "sk_test_51NXI07EcYmsYVNOjkRA66IiQgIDt1ykK42IKEhqFpI7Z8R53ziPcNJ3hlWlqme3YRXtKL8ZmFobslouN1Trr5WdU009OzHQyiM"
	app.Post("/create-payment-intent", func(c *fiber.Ctx) error {
		var req struct {
			DonationAmt int64  `json:"donationAmt"`
			ProductId   string `json:"productId"`
		}

		if err := c.BodyParser(&req); err != nil {
			fmt.Println("Error Parsing Body,: %w", err)
			return c.SendStatus(500)
		}

		orderTotal := int64(req.DonationAmt)
		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(orderTotal),
			Currency: stripe.String(string(stripe.CurrencyUSD)),
			// In the latest version of the API, specifying the `automatic_payment_methods` parameter is optional because Stripe enables its functionality by default.
			AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
				Enabled: stripe.Bool(true),
			},
			Metadata: map[string]string{
				"product":     req.ProductId,
				"donationAmt": strconv.FormatInt(req.DonationAmt, 10),
			},
		}

		pi, err := paymentintent.New(params)
		fmt.Println("pi.New: %w", pi.ClientSecret)

		if err != nil {
			fmt.Println("pi.New: %w", err)
			return c.SendStatus(500)
		}

		return c.JSON(fiber.Map{
			"clientSecret": pi.ClientSecret,
		})
	})

	app.Post("/payment-success", func(c *fiber.Ctx) error {
		var piSuccess struct {
			ID       string `json:"id"`
			Amount   int32  `json:"amount"`
			Email    string `json:"receipt_email"`
			Username string `json:"username"`
		}

		c.BodyParser(&piSuccess)
		fmt.Println(piSuccess)
		pi, err := paymentintent.Get(piSuccess.ID, &stripe.PaymentIntentParams{})
		fmt.Println(pi.Status)

		if pi.Status == stripe.PaymentIntentStatusProcessing || pi.Status == stripe.PaymentIntentStatusSucceeded {
			err = db.Db.CreateUser(c.Context(), sqlc.CreateUserParams{
				Email:           piSuccess.Email,
				Username:        piSuccess.Username,
				PaymentIntent:   piSuccess.ID,
				DonationInCents: piSuccess.Amount,
			})
			var pgErr *pgconn.PgError
			if err != nil {
				if errors.As(err, &pgErr) {
					if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
						db.Db.AddToBalance(c.Context(), sqlc.AddToBalanceParams{Email: piSuccess.Email, AddAmount: piSuccess.Amount})
						fmt.Println("Email Already Exists, DonationAmt Updated")
						return c.JSON(fiber.Map{"status": "fail", "message": "Email already has an account, Donation Amount has been updated"})
					}
				} else {
					fmt.Println(err)
				}
				return c.JSON(fiber.Map{"status": "fail"})

			}
			return c.JSON(fiber.Map{"status": "success"})
		}
		if err != nil {
			fmt.Println(err)
		}
		return c.JSON(fiber.Map{"status": "fail"})

	})
	app.Post("/init-password", func(c *fiber.Ctx) error {
		var initPassword struct {
			Password string `json:"assword"`
			Pi       string `json:"paymentId"`
		}
		c.BodyParser(&initPassword)
		fmt.Println("init: ", initPassword)
		passHash, err := bcrypt.GenerateFromPassword([]byte(initPassword.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Error Hashing Password: ", err)
			return c.JSON(fiber.Map{"message": "error hashing password"})
		}
		id, err := db.Db.InitPass(c.Context(), sqlc.InitPassParams{PaymentIntent: initPassword.Pi, PassHash: passHash})
		if err != nil {
			fmt.Println("Error Initing Password: ", err)
			return c.SendString(`{"message":"Did you already set your password?"}`)
		}
		// fmt.Println(id)
		cookies, err := auth.CreateLoginCookies(id)
		if err != nil {
			fmt.Println("cookie creation error: %w", err)
			return c.JSON(fiber.Map{"error": "Error Signing You in. Please Try Again"})
		}
		for _, cookie := range cookies {
			c.Cookie(cookie)
		}
		return c.JSON(fiber.Map{"go-to": "/super-cool-kids"})

	})

}
