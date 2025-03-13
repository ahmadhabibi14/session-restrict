package service

import (
	"net/http"
	"session-restrict/src/dto/request"
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

func (s *Session) Approve(in request.ReqSessionApprove) (out response.ResSessionApprove, err error) {
	sess := sessions.NewSession()

	key := sess.GenerateKey(in.Role, in.UserId, in.AccessToken)
	session, err := sess.Approve(key)
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		return
	}

	out.Data = &session

	return
}

func (s *Session) Delete(in request.ReqSessionDelete) (out response.ResponseCommon, err error) {
	sess := sessions.NewSession()

	key := sess.GenerateKey(in.Role, in.UserId, in.AccessToken)
	err = sess.DeleteSession(key)
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		return
	}

	return
}
