package sessions

import (
	"errors"
	"fmt"
	"session-restrict/helper"
	"session-restrict/src/lib/database"
	"session-restrict/src/lib/logger"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

const (
	RoleAdmin  = `admin`
	RoleDriver = `driver`
	RoleUser   = `user`
)

type Session struct {
	UserId    uint64    `json:"user_id"`
	Role      string    `json:"role"`
	IpV4      string    `json:"ip_v4"`
	IpV6      string    `json:"ip_v6"`
	UserAgent string    `json:"user_agent"`
	Device    string    `json:"device"`
	OS        string    `json:"os"`
	Approved  bool      `json:"approved"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	Err500FailedSetSession    = errors.New(`failed to set session`)
	Err500FailedGetSession    = errors.New(`failed to get session`)
	Err500FailedGetSessions   = errors.New(`failed to get sessions`)
	Err500FailedDeleteSession = errors.New(`failed to delete session`)
	Err400InvalidToken        = errors.New(`invalid access token`)
)

// session:<role>:<user_id>:<access_token>
func GetKey(role string, userId uint64, accessToken string) string {
	return fmt.Sprintf("session:%v:%d:%s", role, userId, accessToken)
}

func GetDuration(future time.Time) time.Duration {
	return future.Sub(time.Now())
}

func SetSession(session Session, role string, userId uint64, expired time.Duration) (string, error) {
	sessionJson, err := json.Marshal(session)
	if err != nil {
		logger.Log.Error(err, Err500FailedSetSession.Error())
		return ``, Err500FailedSetSession
	}

	token := helper.RandString(20)
	key := GetKey(role, userId, token)

	err = database.ConnRd.Set(key, sessionJson, expired).Err()

	if err != nil {
		logger.Log.Error(err, Err500FailedSetSession.Error())
		return ``, Err500FailedSetSession
	}

	return token, nil
}

func GetSessionByToken(accessToken string) (session Session, err error) {
	keyPattern := fmt.Sprintf("*%s", accessToken)

	var scannedKeys []string

	iter := database.ConnRd.Scan(0, keyPattern, 10).Iterator()
	for iter.Next() {
		scannedKeys = append(scannedKeys, iter.Val())
	}

	if len(scannedKeys) == 0 {
		err = Err400InvalidToken
		return
	}

	for _, k := range scannedKeys {
		var sessString string
		errGet := database.ConnRd.Get(k).Scan(&sessString)
		if errGet != nil {
			logger.Log.Error(errGet, Err500FailedGetSession.Error())

			err = Err500FailedGetSession
			return
		}

		errUnmarshal := json.Unmarshal([]byte(sessString), &session)
		if errUnmarshal != nil {
			logger.Log.Error(errUnmarshal, Err500FailedGetSession.Error())

			err = Err500FailedGetSession
			return
		}

		break
	}

	return
}

type SessionsWithKey struct {
	Session
	Key         string `json:"key,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}

func GetSessionByRoleByUserId(role string, userId uint64) (exist bool, session SessionsWithKey, err error) {
	keyPattern := fmt.Sprintf("session:%s:%d:*", role, userId)

	var scannedKeys []string

	iter := database.ConnRd.Scan(0, keyPattern, 10).Iterator()
	for iter.Next() {
		scannedKeys = append(scannedKeys, iter.Val())
	}

	if len(scannedKeys) > 0 {
		exist = true
	}

	for _, k := range scannedKeys {
		var sessString string
		errGet := database.ConnRd.Get(k).Scan(&sessString)
		if errGet != nil {
			logger.Log.Error(errGet, Err500FailedGetSession.Error())

			err = Err500FailedGetSession
			return
		}

		errUnmarshal := json.Unmarshal([]byte(sessString), &session)
		if errUnmarshal != nil {
			logger.Log.Error(errUnmarshal, Err500FailedGetSession.Error())

			err = Err500FailedGetSession
			return
		}

		session.Key = k
		session.AccessToken = getTokenFromKey(k)

		if !session.Approved {
			break
		}
	}

	return
}

func GetSessionsByRoleByUserId(role string, userId uint64) (sessions []SessionsWithKey, err error) {
	keyPattern := fmt.Sprintf("session:%s:%d:*", role, userId)

	var scannedKeys []string

	iter := database.ConnRd.Scan(0, keyPattern, 10).Iterator()
	for iter.Next() {
		scannedKeys = append(scannedKeys, iter.Val())
	}

	for _, k := range scannedKeys {
		var sessString string
		errGet := database.ConnRd.Get(k).Scan(&sessString)
		if errGet != nil {
			logger.Log.Error(errGet, Err500FailedGetSessions.Error())

			err = Err500FailedGetSessions
			return
		}

		var session SessionsWithKey
		errUnmarshal := json.Unmarshal([]byte(sessString), &session)
		if errUnmarshal != nil {
			logger.Log.Error(errUnmarshal, Err500FailedGetSessions.Error())

			err = Err500FailedGetSessions
			return
		}

		session.Key = k
		session.AccessToken = getTokenFromKey(k)

		sessions = append(sessions, session)
	}

	return
}

func DeleteSession(key string) error {
	err := database.ConnRd.Del(key).Err()
	if err != nil {
		logger.Log.Error(err)

		return Err500FailedDeleteSession
	}

	return nil
}

func getTokenFromKey(key string) string {
	parts := strings.Split(key, ":") // Split by ":"

	if len(parts) >= 4 {
		token := parts[3] // Get the last part
		return token
	}

	return key
}
