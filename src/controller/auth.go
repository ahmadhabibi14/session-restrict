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
	srvAuth *service.Auth
}

func NewAuth(app *fiber.App, srvAuth *service.Auth) {
	handler := &Auth{srvAuth}

	app.Route("/api/auth", func(router fiber.Router) {
		router.Post("/signin", handler.SignIn)
		router.Post("/signup", handler.SignUp)
		router.Post("/signout", handler.SignOut)
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

	in.IpV4 = GetIpV4(c)
	in.IpV6 = GetIpV6(c)
	in.UserAgent = c.Get(fiber.HeaderUserAgent)
	in.OS = GetOS(c)
	in.Device = GetDevice(c)

	out, err := a.srvAuth.SignIn(in)
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
	in, err := helper.ReadBody[request.ReqAuthSignUp](c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.ResponseCommon{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
	}

	out, err := a.srvAuth.SignUp(in)
	if err != nil {
		return c.Status(out.StatusCode).JSON(response.ResponseCommon{
			StatusCode: out.StatusCode,
			Error:      err.Error(),
		})
	}

	out.SetMessage(`Sign Up successfully !`)
	out.SetStatus(http.StatusOK)

	return c.Status(http.StatusOK).JSON(out)
}

func (a *Auth) SignOut(c *fiber.Ctx) error {
	session := getSession(c)

	out, err := a.srvAuth.SignOut(session.UserId, session.AccessToken, session.Role)
	if err != nil {
		return c.Status(out.StatusCode).JSON(response.ResponseCommon{
			StatusCode: out.StatusCode,
			Error:      err.Error(),
		})
	}

	RemoveAuthCookie(c)

	out.SetMessage(`Sign Out Successful !`)
	out.SetStatus(http.StatusOK)

	return c.Status(http.StatusOK).JSON(out)
}
