package controller

import (
	"net/http"
	"session-restrict/src/service"

	"github.com/gofiber/fiber/v2"
)

type Auth struct {
	auth *service.Auth
}

func NewAuth(app *fiber.App, auth *service.Auth) {
	handler := &Auth{auth}

	app.Route("/api/auth", func(router fiber.Router) {
		router.Post("/signin", handler.SignIn)
		router.Post("/signup", handler.SignUp)
	})
}

func (a *Auth) SignIn(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

func (a *Auth) SignUp(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}
