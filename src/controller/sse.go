package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetSSEPayload(event string, data string) string {
	return fmt.Sprintf("event: %v\ndata: %s\n\n", event, data)
}

func SetSSEHeaders(c *fiber.Ctx) {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	c.Set("Content-Encoding", "none")
}
