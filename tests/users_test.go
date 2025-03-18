package tests

import (
	"session-restrict/helper/converter"
	"session-restrict/src/lib/database"
	"session-restrict/src/repo/users"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	assert.NotNil(t, database.ConnPg, `postgres connection is empty`)
	assert.NotNil(t, database.ConnRd, `redis connection is empty`)

	t.Run(`insert`, func(t *testing.T) {
		userFullName := `Kaito Kuroba`
		userEmail := `kaito1412@example.com`
		userPassword := `kaito#1412`

		user := users.NewUser()
		user.FullName = userFullName
		user.Email = userEmail
		user.Password = userPassword
		user.Role = users.RoleUser

		assert.NoError(t, user.Insert(), `failed to insert user`)

		t.Log(`Inserted User : `, converter.AnyToJsonPretty(user))

		t.Run(`findByEmail`, func(t *testing.T) {
			userFindByEmail := users.NewUser()
			userFindByEmail.Email = userEmail

			assert.NoError(t, userFindByEmail.FindByEmail(), `failed to find user by email`)

			t.Log(`Found User By Email : `, converter.AnyToJsonPretty(userFindByEmail))
		})
	})
}
