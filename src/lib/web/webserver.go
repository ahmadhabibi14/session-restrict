package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
)

type WebServer struct {
	fiber.Config
}

func NewWebserver() *fiber.App {
	engine := django.New("./src/views", ".django")

	return fiber.New(fiber.Config{
		AppName:                 "Dummy Session Restriction",
		Views:                   engine,
		Prefork:                 false,
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		EnableTrustedProxyCheck: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			var code int = fiber.StatusNotFound
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			if c.Method() == fiber.MethodGet {
				if code == fiber.StatusNotFound {
					return c.Render("error", fiber.Map{
						`Title`:       fmt.Sprintf("%d - %s", fiber.StatusNotFound, `Page Not Found`),
						`Description`: "Cannot find the page you are looking for",
					})
				}

				return c.Render("error", fiber.Map{
					`Title`:       fmt.Sprintf("%d - %s", code, http.StatusText(code)),
					`Description`: err.Error(),
				})
			}

			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return c.Status(code).JSON(fiber.Map{
				"errors": e.Error(),
			})
		},
		Immutable: true,
		BodyLimit: 40 * 1024 * 1024,
	})
}
