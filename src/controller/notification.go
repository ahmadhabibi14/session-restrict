package controller

import (
	"bufio"
	"session-restrict/helper/converter"
	"session-restrict/src/lib/database"
	"session-restrict/src/lib/logger"
	"session-restrict/src/repo/sessions"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type Notification struct {
}

func NewNotification(app *fiber.App) {
	handler := &Notification{}

	app.Route("/api/notification", func(router fiber.Router) {
		router.Get("/user", mustLoggedInAjax, handler.ByUserId)
	})
}

func (n *Notification) ByUserId(c *fiber.Ctx) error {
	SetSSEHeaders(c)

	session := getSession(c)

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		channel := sessions.GetChannelUserNotif(session.UserId)
		pubsub := database.ConnRd.Subscribe(channel)
		defer pubsub.Close()

		for {
			msg, err := pubsub.ReceiveMessage()
			if err != nil {
				logger.Log.Error(err)
				continue
			}

			var out fiber.Map
			err = json.Unmarshal([]byte(msg.Payload), &out)
			if err != nil {
				logger.Log.Error(err)
				continue
			}

			event := converter.AnyToString(out[`event`])
			var dataBytes []byte

			switch event {
			case sessions.EventNewSession:
				var notifSession sessions.NotificationNewSession
				err = json.Unmarshal([]byte(msg.Payload), &notifSession)
				if err != nil {
					logger.Log.Error(err)
					continue
				}

				dataBytes, err = json.Marshal(notifSession.Data)
				if err != nil {
					logger.Log.Error(err)
					continue
				}
			default:
				continue
			}

			payload := GetSSEPayload(sessions.EventNewSession, string(dataBytes))
			_, err = w.WriteString(payload)
			if err != nil {
				logger.Log.Error(err, `failed to write data`)
				return
			}

			err = w.Flush()
			if err != nil {
				logger.Log.Error(err, `failed to flush data`)
				return
			}
		}
	}))

	return nil
}
