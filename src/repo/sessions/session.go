package sessions

import (
	"errors"
	"fmt"
	"session-restrict/helper"
	"session-restrict/src/lib/database"
	"session-restrict/src/lib/logger"
	"time"

	"github.com/goccy/go-json"
)

const (
	RoleAdmin  = `admin`
	RoleDriver = `driver`
	RoleUser   = `user`
)

type Session struct {
	UserId    uint64 `json:"user_id"`
	Role      string `json:"role"`
	IpV4      string `json:"ip_v4"`
	IpV6      string `json:"ip_v6"`
	UserAgent string `json:"user_agent"`
	Device    string `json:"device"`
	OS        string `json:"os"`
}

var (
	Err500FailedSetSession = errors.New(`failed to set session`)
	Err400InvalidToken     = errors.New(`invalid access token`)
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

func GetSessionByToken(accessToken string) (Session, error) {
	var sess Session
	var sessStr string

	keyPattern := fmt.Sprintf("session:*:*:%s", accessToken)

	err := database.ConnRd.Get(keyPattern).Scan(&sessStr)
	if !(err == nil || sessStr == "") {
		logger.Log.Error(err, Err400InvalidToken.Error())

		return sess, Err400InvalidToken
	}

	err = json.Unmarshal([]byte(sessStr), &sess)
	if err != nil {
		logger.Log.Error(err, Err500FailedSetSession.Error())
		return sess, Err500FailedSetSession
	}

	return sess, nil
}

// func GetSessionsByRoleByUserId(role string, userId uint64) ([]Session, error) {
// 	var sess Session
// 	var sessStr string

// 	keyPattern := fmt.Sprintf("session:%s:%d:*", role, userId)

// 	iter := database.ConnRd.Scan()

// 	return sess, nil
// }
