package sessions

import (
	"errors"
	"fmt"
	"session-restrict/src/lib/database"
	"session-restrict/src/lib/logger"

	"github.com/goccy/go-json"
)

const (
	EventNewSession         = `new_session`
	EventNewSessionApproved = `new_session_approved`
	EventNewSessionDeleted  = `new_session_deleted`
)

func GetChannelUserNotif(userId uint64) string {
	return fmt.Sprintf("user.notification.%v", userId)
}

type NotificationNewSession struct {
	Event string  `json:"event"`
	Data  Session `json:"data"`
}

type NotificationNewSessionApproved struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type NotificationNewSessionDeleted struct {
	Event string `json:"event"`
	Data  string `json:"data"`
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

func PublishNewSessionApproved(data string, userId uint64) error {
	channel := GetChannelUserNotif(userId)
	in := NotificationNewSessionApproved{
		Event: EventNewSessionApproved,
		Data:  data,
	}
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

func PublishNewSessionDeleted(data string, userId uint64) error {
	channel := GetChannelUserNotif(userId)
	in := NotificationNewSessionDeleted{
		Event: EventNewSessionDeleted,
		Data:  data,
	}
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
