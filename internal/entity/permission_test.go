package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type permissionTestContext struct {
	sut      *Permission
	name     string
	codename string
}

func (tc *permissionTestContext) setUp() {
	tc.name = "Criar Todo"
	tc.codename = "todo.add"
	tc.sut = NewPermission(tc.name, tc.codename)
}

func (tc *permissionTestContext) tearDown() {
	tc.sut = nil
}

func TestNewPermission(t *testing.T) {
	tc := &permissionTestContext{}
	t.Run("Should create a new permission with correct params", func(t *testing.T) {
		tc.setUp()
		defer tc.tearDown()
		assert.NotNil(t, tc.sut)
		assert.NotEmpty(t, tc.sut.ID)
		assert.Equal(t, tc.name, tc.sut.Name)
		assert.Equal(t, tc.codename, tc.sut.Codename)
	})
}
