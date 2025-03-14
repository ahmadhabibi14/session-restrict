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
	AccessToken string    `json:"access_token"`
	UserId      uint64    `json:"user_id"`
	Role        string    `json:"role"`
	IpV4        string    `json:"ip_v4"`
	IpV6        string    `json:"ip_v6"`
	UserAgent   string    `json:"user_agent"`
	Device      string    `json:"device"`
	OS          string    `json:"os"`
	Approved    bool      `json:"approved"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

func NewSession() *Session {
	return &Session{}
}

var (
	Err500FailedSetSession    = errors.New(`failed to set session`)
	Err500FailedGetSession    = errors.New(`failed to get session`)
	Err500FailedGetSessions   = errors.New(`failed to get sessions`)
	Err500FailedDeleteSession = errors.New(`failed to delete session`)
	Err500FailedUpdateSession = errors.New(`failed to update session`)
	Err400SessionNotFound     = errors.New(`session not found`)
	Err400InvalidToken        = errors.New(`invalid access token`)
)

func (s *Session) GenerateToken() string {
	return helper.RandString(30)
}

func (s *Session) GenerateKey(role string, userId uint64, accessToken string) string {
	return fmt.Sprintf("session:%v:%d:%s", role, userId, accessToken)
}

func (s *Session) GenerateDuration(future time.Time) time.Duration {
	return future.Sub(time.Now())
}

func (s *Session) GetSession(key string) (Session, error) {
	var sessString string
	var sess Session
	var err error

	err = database.ConnRd.Get(key).Scan(&sessString)
	if err != nil || sessString == "" {
		return sess, Err400SessionNotFound
	}

	err = json.Unmarshal([]byte(sessString), s)
	if err != nil {
		logger.Log.Error(err)
		return sess, Err500FailedGetSession
	}

	return sess, nil
}

func (s *Session) SetSession(expired time.Duration) error {
	sessionJson, err := json.Marshal(s)
	if err != nil {
		logger.Log.Error(err, Err500FailedSetSession.Error())
		return Err500FailedSetSession
	}

	key := s.GenerateKey(s.Role, s.UserId, s.AccessToken)

	err = database.ConnRd.Set(key, sessionJson, expired).Err()
	if err != nil {
		logger.Log.Error(err, Err500FailedSetSession.Error())
		return Err500FailedSetSession
	}

	return nil
}

func (s *Session) Approve(key string) (Session, error) {
	var sessString string
	var sess Session
	var err error

	err = database.ConnRd.Get(key).Scan(&sessString)
	if err != nil || sessString == "" {
		return sess, Err400SessionNotFound
	}

	err = json.Unmarshal([]byte(sessString), &sess)
	if err != nil {
		logger.Log.Error(err)
		return sess, Err500FailedGetSession
	}

	sess.Approved = true
	sess.UpdatedAt = time.Now()

	sessJson, err := json.Marshal(sess)
	if err != nil {
		logger.Log.Error(err)
		return sess, Err500FailedGetSession
	}

	ttl, err := database.ConnRd.TTL(key).Result()
	if err != nil || ttl <= 0 {
		logger.Log.Error(err)
		return sess, Err500FailedUpdateSession
	}

	err = database.ConnRd.Set(key, sessJson, ttl).Err()
	if err != nil {
		logger.Log.Error(err)
		return sess, Err500FailedUpdateSession
	}

	return sess, nil
}

func (s *Session) GetSessionByToken() (session Session, err error) {
	keyPattern := fmt.Sprintf("*%s", s.AccessToken)

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

func (s *Session) GetSessionByRoleByUserId() (session Session, isExist bool, err error) {
	keyPattern := fmt.Sprintf("session:%s:%d:*", s.Role, s.UserId)

	var scannedKeys []string

	iter := database.ConnRd.Scan(0, keyPattern, 10).Iterator()
	for iter.Next() {
		scannedKeys = append(scannedKeys, iter.Val())
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

		if session.Approved {
			isExist = true
			break
		}
	}

	return
}

func (s *Session) GetSessionsByUser() (sessions []Session, err error) {
	keyPattern := fmt.Sprintf("session:%s:%d:*", s.Role, s.UserId)

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

		var session Session
		errUnmarshal := json.Unmarshal([]byte(sessString), &session)
		if errUnmarshal != nil {
			logger.Log.Error(errUnmarshal, Err500FailedGetSessions.Error())

			err = Err500FailedGetSessions
			return
		}

		sessions = append(sessions, session)
	}

	return
}

func (s *Session) DeleteSession(key string) error {
	err := database.ConnRd.Del(key).Err()
	if err != nil {
		logger.Log.Error(err)

		return Err500FailedDeleteSession
	}

	return nil
}
