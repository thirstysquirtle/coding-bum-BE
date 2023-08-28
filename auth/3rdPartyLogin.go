package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func Set3rdPartyLogin(app *fiber.App) {

	conf := &oauth2.Config{
		ClientID:     "914923693614-va28h0qi3cms03tghvnibmdc2hbd517j.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:3000/ass",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	app.Get("/gae", func(c *fiber.Ctx) error {
		url := conf.AuthCodeURL("statasde", oauth2.AccessTypeOnline)
		return c.Redirect(url)
	})

	app.Use("/ass", func(c *fiber.Ctx) error {
		// fmt.Println(c.Queries())
		// fmt.Println(c.FormValue("code"), c.FormValue("state"))
		token, _ := conf.Exchange(c.Context(), c.FormValue("code"))
		fmt.Println(token)
		// if er != nil {
		// 	fmt.Println(er)
		// }
		return c.JSON(*token)
	})

}
