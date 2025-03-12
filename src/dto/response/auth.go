package response

import (
	"session-restrict/src/repo/users"
	"time"
)

type ResAuthSignIn struct {
	ResponseCommon
	User        *users.User `json:"user"`
	AccessToken string      `json:"access_token"`
	ExpiredAt   time.Time   `json:"expired_at"`
} // @name ResAuthSignIn

type ResAuthSignUp struct {
	ResponseCommon
	User *users.User `json:"user"`
}
