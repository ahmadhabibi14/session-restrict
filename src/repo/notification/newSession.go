package notification

import (
	"errors"
	"fmt"
	"session-restrict/src/lib/database"
	"session-restrict/src/lib/logger"
	"time"

	"github.com/goccy/go-json"
)

const (
	EventNewSession     = `new_session`
	EventSessionRemoved = `session_removed`
)

func GetChannelUserNotif(userId uint64) string {
	return fmt.Sprintf("user.notification.%d", userId)
}

type NewSessionData struct {
	UserId      uint64    `json:"user_id"`
	Role        string    `json:"role"`
	AccessToken string    `json:"access_token"`
	Timestamp   time.Time `json:"timestamp"`
}

type NewSession struct {
	Event string         `json:"event"`
	Data  NewSessionData `json:"data"`
}

var ErrPulishNewSession = errors.New(`failed to send notification`)

func PublishNewSession(in NewSession, userId uint64) error {
	channel := GetChannelUserNotif(userId)
	dataByte, err := json.Marshal(in)
	if err != nil {
		logger.Log.Error(err)

		return ErrPulishNewSession
	}

	err = database.ConnRd.Publish(channel, dataByte).Err()
	if err != nil {
		logger.Log.Error(err)

		return ErrPulishNewSession
	}

	return nil
}
