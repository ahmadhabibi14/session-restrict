package helper

import (
	"errors"
	"session-restrict/src/app/logger"

	"github.com/gofiber/fiber/v2"
)

func ReadBody[T any | struct{}](c *fiber.Ctx) (out T, err error) {
	err = c.BodyParser(&out)
	if err != nil {
		logger.Log.Error(err, `failed to parse request body`)
		err = errors.New(`invalid payload, please check your request and try again`)
		return
	}

	err = ValidateStruct(out)
	if err != nil {
		return
	}

	return
}

func ReadQuery[T any | struct{}](c *fiber.Ctx) (out T, err error) {
	if err = c.QueryParser(&out); err != nil {
		logger.Log.Error(err, `failed to parse request query`)
		err = errors.New(`invalid payload, please check your request and try again`)
		return
	}

	err = ValidateStruct(out)
	if err != nil {
		return
	}

	return
}
