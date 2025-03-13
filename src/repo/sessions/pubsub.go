package sessions

import (
	"errors"
	"fmt"
	"session-restrict/src/lib/database"
	"session-restrict/src/lib/logger"

	"github.com/goccy/go-json"
)

const (
	EventNewSession     = `new_session`
	EventSessionRemoved = `session_removed`
)

func GetChannelUserNotif(userId uint64) string {
	return fmt.Sprintf("user.notification.%d", userId)
}

type NotificationNewSession struct {
	Event string  `json:"event"`
	Data  Session `json:"data"`
}

var ErrPublishNewSession = errors.New(`failed to send notification`)

func PublishNewSession(in NotificationNewSession, userId uint64) error {
	channel := GetChannelUserNotif(userId)
	dataByte, err := json.Marshal(in)
	if err != nil {
		logger.Log.Error(err)

		return ErrPublishNewSession
	}

	err = database.ConnRd.Publish(channel, dataByte).Err()
	if err != nil {
		logger.Log.Error(err)

		return ErrPublishNewSession
	}

	return nil
}
