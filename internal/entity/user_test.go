package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("John Doe", "j@j.com", "123456")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "j@j.com", user.Email)
}

func TestUserWhenNameIsRequired(t *testing.T) {
	u, err := NewUser("", "j@j.com", "123")
	assert.Nil(t, u)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestUserWhenEmailIsRequired(t *testing.T) {
	u, err := NewUser("teste", "", "123")
	assert.Nil(t, u)
	assert.Equal(t, ErrEmailIsRequired, err)
}

func TestUserWhenPasswordIsRequired(t *testing.T) {
	u, err := NewUser("teste", "j@j.com", "")
	assert.Nil(t, u)
	assert.Equal(t, ErrPasswordIsRequired, err)
}

func TestUser_ValidatePassword(t *testing.T) {
	user, err := NewUser("John Doe", "j@j.com", "123456")
	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword("123456"))
	assert.False(t, user.ValidatePassword("1234567"))
	assert.NotEqual(t, "123456", user.Password)
}
