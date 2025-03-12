package users

import (
	"errors"
	"session-restrict/src/lib/database"
	"session-restrict/src/lib/logger"
	"time"

	"github.com/lib/pq"
)

const TableUsers = `users`

const (
	RoleAdmin  = `admin`
	RoleUser   = `user`
	RoleDriver = `driver`
)

func IsValidRole(role string) bool {
	switch role {
	case RoleAdmin, RoleUser, RoleDriver:
		return true
	default:
		return false
	}
}

type User struct {
	Id        uint64    `db:"id" json:"id"`
	FullName  string    `db:"full_name" json:"full_name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"`
	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted bool      `db:"is_deleted" json:"is_deleted"`
}

func NewUser() *User {
	return &User{}
}

var (
	Err400UserInsertEmailExist  = errors.New(`email already exist`)
	Err400UserInsertInvalidRole = errors.New(`invalid user role`)
	Err500UserInsertFailed      = errors.New(`error while creating user`)
)

func (u *User) Insert() error {
	if !IsValidRole(u.Role) {
		return Err400UserInsertInvalidRole
	}

	query := `INSERT INTO users  (full_name, email, password, role)
	VALUES ($1, $2, $3, $4) RETURNING *`

	if err := database.ConnPg.QueryRowx(query,
		u.FullName, u.Email, u.Password, u.Role,
	).StructScan(u); err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == pq.ErrorCode(`23505`) {
			logger.Log.Error(pgErr)
			return Err400UserInsertEmailExist
		}

		logger.Log.Error(err, Err500UserInsertFailed.Error())
		return Err500UserInsertFailed
	}

	return nil
}
