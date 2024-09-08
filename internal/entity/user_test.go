package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeUser() (*User, error) {
	return NewUser("Carlos", "t@t.com", "123456")
}

func TestNewUser(t *testing.T) {
	user, err := makeUser()
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "Carlos", user.Name)
	assert.Equal(t, "t@t.com", user.Email)
}

func TestUser_ValidatePassword(t *testing.T) {
	user, err := makeUser()
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.True(t, user.ValidatePassword("123456"))
	assert.False(t, user.ValidatePassword("1234567"))
	assert.NotEqual(t, "123456", user.Password)
}
