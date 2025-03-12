package sessions

import (
	"encoding/json"
	"errors"
	"session-restrict/helper"
	"session-restrict/src/lib/database"
	"session-restrict/src/lib/logger"
	"strings"
	"time"
)

// Session key in Redis:
// session:<role>:<refresh-token>

type Role string

const (
	RoleAdmin  Role = `admin`
	RoleDriver Role = `driver`
	RoleUser   Role = `user`
)

type Session struct {
	UserId    uint64 `json:"user_id"`
	Role      string `json:"role"`
	IpAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Device    string `json:"device"`
	OS        string `json:"os"`
}

func GenerateKey(role Role, token string) string {
	switch role {
	case RoleAdmin, RoleDriver, RoleUser:
		break
	default:
		role = RoleUser
	}

	return `session:` + string(role) + `:` + token
}

func GetKey(role Role, token string) string {
	switch role {
	case RoleAdmin, RoleDriver, RoleUser:
		break
	default:
		role = RoleUser
	}
	return `session:` + string(role) + `:` + token
}

var Err500SetSessionAdminFailed = errors.New(`failed to set session admin`)

func SetSession(session Session, role Role, expired time.Duration) (string, error) {
	sessionJson, err := json.Marshal(session)
	if err != nil {
		logger.Log.Error(err, Err500SetSessionAdminFailed.Error())
		return ``, Err500SetSessionAdminFailed
	}

	token := helper.RandString(20)
	key := GenerateKey(RoleAdmin, token)
	err = database.ConnRd.Set(key, sessionJson, expired).Err()
	if err != nil {
		logger.Log.Error(err, Err500SetSessionAdminFailed.Error())
		return ``, Err500SetSessionAdminFailed
	}

	return token, nil
}

var (
	Err500GetSessionAdminFailed   = errors.New(`failed to get session admin`)
	Err400GetSessionAdminNotFound = errors.New(`invalid refresh token`)
)

func GetSession(token string) (session Session, err error) {
	key := GetKey(RoleAdmin, token)

	var sessionString string
	errGet := database.ConnRd.Get(key).Scan(&sessionString)
	if errGet != nil {
		logger.Log.Error(errGet, Err400GetSessionAdminNotFound.Error())
		return Session{}, Err400GetSessionAdminNotFound
	}

	if sessionString == `` {
		return Session{}, Err400GetSessionAdminNotFound
	}

	err = json.Unmarshal([]byte(sessionString), &session)
	if err != nil {
		logger.Log.Error(err, Err500GetSessionAdminFailed.Error())
		return Session{}, Err500GetSessionAdminFailed
	}

	return
}

var (
	Err500GetSessionsAdminFailed = errors.New(`failed to get sessions admin`)
)

type SessionWithKey struct {
	Session
	Key string `json:"-"`
}

func GetSessionsAdmin() ([]SessionWithKey, error) {
	var cursor uint64
	var keys []string
	var pattern = `session:` + string(RoleAdmin) + `:*`
	var err error

	var sessions = []SessionWithKey{}

	for {
		var scannedKeys []string
		scannedKeys, cursor, err = database.ConnRd.Scan(cursor, pattern, 10).Result()
		if err != nil {
			logger.Log.Error(err, Err500GetSessionsAdminFailed.Error())
			return sessions, Err500GetSessionsAdminFailed
		}

		keys = append(keys, scannedKeys...)

		if cursor == 0 {
			break
		}
	}

	for _, key := range keys {
		var sessionString string
		errGet := database.ConnRd.Get(key).Scan(&sessionString)
		if errGet != nil {
			logger.Log.Error(errGet, Err500GetSessionsAdminFailed.Error())
			continue
		}

		if sessionString == `` {
			continue
		}

		var session SessionWithKey
		err = json.Unmarshal([]byte(sessionString), &session)
		if err != nil {
			logger.Log.Error(err, Err500GetSessionsAdminFailed.Error())
			continue
		}

		session.Key = key

		sessions = append(sessions, session)
	}

	return sessions, nil
}

var (
	Err500DeleteAdminSessionFailed       = errors.New(`failed to delete session admin`)
	Err400DeleteAdminSessionInvalidToken = errors.New(`invalid refresh token`)
)

func DeleteAdminSession(token string) error {
	if !strings.Contains(token, string(RoleAdmin)) {
		token = GetKey(RoleAdmin, token)
	}

	var sessionString string
	err := database.ConnRd.Get(token).Scan(&sessionString)
	if err != nil || sessionString == "" {
		return Err400DeleteAdminSessionInvalidToken
	}

	err = database.ConnRd.Del(token).Err()
	if err != nil {
		logger.Log.Error(err, Err500DeleteAdminSessionFailed.Error())
		return Err500DeleteAdminSessionFailed
	}

	return nil
}
