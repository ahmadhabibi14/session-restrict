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

func (s *Session) Approve(in request.ReqSessionApprove, userId uint64) (out response.ResSessionApprove, err error) {
	sess := sessions.NewSession()

	key := sess.GenerateKey(in.Role, userId, in.AccessToken)
	session, err := sess.Approve(key)
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		return
	}

	err = sessions.PublishNewSessionApproved(sessions.NotificationNewSessionApproved{
		Event: sessions.EventNewSessionApproved,
		Data:  *sess,
	}, in.UserId)
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		return
	}

	out.Data = &session

	return
}

func (s *Session) Delete(in request.ReqSessionDelete, userId uint64) (out response.ResponseCommon, err error) {
	sess := sessions.NewSession()

	key := sess.GenerateKey(in.Role, userId, in.AccessToken)
	session, err := sess.GetSession(key)

	err = sess.DeleteSession(key)
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		return
	}

	err = sessions.PublishNewSessionDeleted(sessions.NotificationNewSessionDeleted{
		Event: sessions.EventNewSessionDeleted,
		Data:  session,
	}, in.UserId)
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		return
	}

	return
}
