package entity

import (
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

type userTestContext struct {
	sut      *User
	err      error
	name     string
	email    string
	password string
}

func (tc *userTestContext) setUp() {
	tc.name = faker.Name()
	tc.email = faker.Email()
	tc.password = faker.Password()
	tc.sut, tc.err = NewUser(tc.name, tc.email, tc.password)
}

func (tc *userTestContext) tearDown() {
	tc.sut = nil
	tc.err = nil
	tc.name = ""
	tc.email = ""
	tc.password = ""
}

func TestNewUser(t *testing.T) {
	tc := &userTestContext{}
	t.Run("Should create a new user with correct params", func(t *testing.T) {
		tc.setUp()
		defer tc.tearDown()
		assert.Nil(t, tc.err)
		assert.NotNil(t, tc.sut)
		assert.NotEmpty(t, tc.sut.ID)
		assert.NotEmpty(t, tc.sut.Password)
		assert.Equal(t, tc.name, tc.sut.Name)
		assert.Equal(t, tc.email, tc.sut.Email)
	})

	t.Run("Should create a new user with valid password", func(t *testing.T) {
		tc.setUp()
		defer tc.tearDown()
		assert.Nil(t, tc.err)
		assert.NotNil(t, tc.sut)
		assert.True(t, tc.sut.ValidatePassword(tc.password))
		assert.False(t, tc.sut.ValidatePassword("123"))
		assert.NotEqual(t, tc.password, tc.sut.Password)
	})
}
