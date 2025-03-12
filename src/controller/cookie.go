package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetAuthCookie(c *fiber.Ctx, tokenString string, expiredAt time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     `auth`,
		Value:    tokenString,
		Expires:  expiredAt,
		SameSite: "Lax",
		Secure:   false,
		HTTPOnly: true,
		Path:     `/`,
	})
}

func RemoveAuthCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:    `auth`,
		Value:   "",
		Path:    `/`,
		Expires: time.Date(-1, 0, 0, 0, 0, 0, 0, time.Local),
	})
}
