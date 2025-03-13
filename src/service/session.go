package service

import (
	"net/http"
	"session-restrict/src/dto/response"
	"session-restrict/src/repo/sessions"
)

type Session struct {
}

func NewSession() *Session {
	return &Session{}
}

func (s *Session) GetSessions(userId uint64, role string) (out response.ResGetSessionsByUserId, err error) {
	sess := sessions.NewSession()
	sess.UserId = userId
	sess.Role = role

	resp, err := sess.GetSessionsByUser()
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		return
	}

	out.Data = &resp

	return
}
