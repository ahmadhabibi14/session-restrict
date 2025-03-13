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

func (s *Session) GetSessionsByUserId(userId uint64, role string) (out response.ResGetSessionsByUserId, err error) {
	resp, err := sessions.GetSessionsByRoleByUserId(role, userId)
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		return
	}

	out.Data = &resp
	return
}
