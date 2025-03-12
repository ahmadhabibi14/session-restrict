package controller

import (
	"net/http"
	"session-restrict/helper"
	"session-restrict/src/dto/request"
	"session-restrict/src/dto/response"
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
	in, err := helper.ReadBody[request.ReqAuthSignIn](c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.ResponseCommon{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
	}

	out, err := a.auth.SignIn(in)
	if err != nil {
		return c.Status(out.StatusCode).JSON(response.ResponseCommon{
			StatusCode: out.StatusCode,
			Error:      err.Error(),
		})
	}

	SetAuthCookie(c, out.AccessToken, out.ExpiredAt)

	out.SetMessage(`Sign In successfully !`)
	out.SetStatus(http.StatusOK)

	return c.Status(http.StatusOK).JSON(out)
}

func (a *Auth) SignUp(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}
