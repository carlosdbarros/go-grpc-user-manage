package entity

import (
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	println("#################### STARTING TESTS ####################")
	println("Before all tests")
	code := m.Run()
	println("After all tests")
	println("#################### FINISHED TESTS ####################")
	os.Exit(code)
}

func setupTest() func() {
	setup()
	return teardown
}

var (
	user     *User
	err      error
	name     = ""
	email    = ""
	password = ""
)

func setup() {
	println("Before each test")
	name = faker.Name()
	email = faker.Email()
	password = faker.Password()
	user, err = NewUser(name, email, password)
	fmt.Printf("User created: %v\n", user)

}

func teardown() {
	println("After each test")
	user = nil
	err = nil
	name = ""
	email = ""
	password = ""
	println("--------------------")
}

func TestNewUser(t *testing.T) {
	defer setupTest()()
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, email, user.Email)
}

func TestUser_ValidatePassword(t *testing.T) {
	defer setupTest()()
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.True(t, user.ValidatePassword(password))
	assert.False(t, user.ValidatePassword("123"))
	assert.NotEqual(t, password, user.Password)
}
