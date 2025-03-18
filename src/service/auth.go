package service

import (
	"errors"
	"net/http"
	"os"
	"session-restrict/helper"
	"session-restrict/src/dto/request"
	"session-restrict/src/dto/response"
	"session-restrict/src/integration/mailer"
	"session-restrict/src/lib/logger"
	"session-restrict/src/repo/sessions"
	"session-restrict/src/repo/users"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) SignIn(in request.ReqAuthSignIn) (out response.ResAuthSignIn, err error) {
	usr := users.NewUser()
	usr.Email = in.Email

	err = usr.FindByEmail()
	if err != nil {
		switch err {
		case users.Err400UserFindByEmailNotFound:
			out.SetStatus(http.StatusUnauthorized)
			return
		default:
			out.SetStatus(http.StatusInternalServerError)
			return
		}
	}

	if !helper.IsValidPassword(in.Password, usr.Password) {
		out.SetStatus(http.StatusBadRequest)
		err = errors.New(`invalid password`)
		return
	}

	sess := sessions.NewSession()
	sess.Role = usr.Role
	sess.UserId = usr.Id
	_, isExist, err := sess.GetSessionByRoleByUserId()
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		return
	}

	future := time.Now().AddDate(0, 2, 0)
	duration := sess.GenerateDuration(future)
	accessToken := sess.GenerateToken()

	sess.AccessToken = accessToken
	sess.UserId = usr.Id
	sess.Role = usr.Role
	sess.IpV4 = in.IpV4
	sess.IpV6 = in.IpV6
	sess.UserAgent = in.UserAgent
	sess.Device = in.Device
	sess.OS = in.OS
	sess.Approved = !isExist
	sess.CreatedAt = time.Now()
	sess.UpdatedAt = time.Now()
	sess.ExpiredAt = future

	err = sess.SetSession(duration)
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		return
	}

	if isExist {
		notifData := sessions.NotificationNewSession{
			Event: sessions.EventNewSession,
			Data:  *sess,
		}

		err = sessions.PublishNewSession(notifData, usr.Id)
		if err != nil {
			out.SetStatus(http.StatusInternalServerError)
			return
		}

		resetPassLink := os.Getenv("WEB_DOMAIN") + "/reset-password"
		emailTitle := "⚠️ Peringatan: Aktivitas Login Tidak Dikenal"
		emailContent := mailer.HtmlOtpNewSessionLoggedIn(
			emailTitle, resetPassLink, usr.FullName, time.Now().Format(time.DateTime),
			(in.Device + `, ` + in.OS), in.IpV4, accessToken,
		)

		mailService, err := mailer.NewMailer()
		if err != nil {
			logger.Log.Error(err)
		}

		if err = mailService.SendMailHTML(
			[]string{usr.Email},
			[]string{},
			emailTitle, emailContent,
		); err != nil {
			logger.Log.Error(err)
		}
	}

	out.AccessToken = accessToken
	out.ExpiredAt = future
	out.User = usr.Sanitize()

	return
}

func (a *Auth) SignUp(in request.ReqAuthSignUp) (out response.ResAuthSignUp, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		err = errors.New(`failed to generate hash password`)
		return
	}

	usr := users.NewUser()
	usr.Email = in.Email
	usr.FullName = in.FullName
	usr.Role = in.Role
	usr.Password = string(hashedPassword)

	err = usr.Insert()
	if err != nil {
		switch err {
		case users.Err400UserInsertEmailExist, users.Err400UserInsertInvalidRole:
			out.SetStatus(http.StatusBadRequest)
			return
		default:
			out.SetStatus(http.StatusInternalServerError)
			return
		}
	}

	out.User = usr.Sanitize()

	return
}

func (a *Auth) SignOut(userId uint64, accessToken, role string) (out response.ResponseCommon, err error) {
	sess := sessions.NewSession()
	sess.UserId = userId
	sess.AccessToken = accessToken
	sess.Role = role

	key := sess.GenerateKey(role, userId, accessToken)
	err = sess.DeleteSession(key)
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		return
	}

	notifData := sessions.NotificationNewSessionDeleted{
		Event: sessions.EventNewSessionDeleted,
		Data:  *sess,
	}

	err = sessions.PublishNewSessionDeleted(notifData, userId)
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		return
	}

	return
}
