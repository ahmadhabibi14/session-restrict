package service

import (
	"errors"
	"net/http"
	"session-restrict/helper"
	"session-restrict/src/dto/request"
	"session-restrict/src/dto/response"
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

	future := time.Now().AddDate(0, 2, 0)
	sessDuration := sessions.GetDuration(future)
	sessToken, err := sessions.SetSession(
		sessions.Session{
			UserId:    usr.Id,
			Role:      usr.Role,
			IpV4:      in.IpV4,
			IpV6:      in.IpV6,
			UserAgent: in.UserAgent,
			Device:    in.Device,
			OS:        in.OS,
		}, usr.Role, usr.Id, sessDuration,
	)
	if err != nil {
		out.SetStatus(http.StatusInternalServerError)
		return
	}

	out.ExpiredAt = future
	out.AccessToken = sessToken
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
