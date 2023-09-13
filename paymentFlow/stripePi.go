package paymentFlow

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/paymentintent"
)

func SetupPaymentFlow(app *fiber.App) {
	// This is your test secret API key.
	stripe.Key = "sk_test_51NXI07EcYmsYVNOjkRA66IiQgIDt1ykK42IKEhqFpI7Z8R53ziPcNJ3hlWlqme3YRXtKL8ZmFobslouN1Trr5WdU009OzHQyiM"
	app.Get("/account-creation-success", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{"email": c.Query("email"), "username": c.Query("username"), "pi": c.Query("payment_intent")})
	})

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
			ID     string `json:"id"`
			Amount int    `json:"amount"`
		}

		c.BodyParser(&piSuccess)
		fmt.Println(piSuccess)
		pi, err := paymentintent.Get(piSuccess.ID, &stripe.PaymentIntentParams{})
		fmt.Println(pi.Status)

		if pi.Status == stripe.PaymentIntentStatusProcessing || pi.Status == stripe.PaymentIntentStatusSucceeded {
			return c.JSON(fiber.Map{"status": "sucess"})
		}
		if err != nil {
			fmt.Println(err)
		}
		return c.JSON(fiber.Map{"status": "fail"})

	})
	// app.Post("/")

}
