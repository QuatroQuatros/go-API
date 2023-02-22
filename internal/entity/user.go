package entity

import (
	"errors"

	"github.com/QuatroQuatros/go-API/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailIsRequired    = errors.New("email is required")
	ErrPasswordIsRequired = errors.New("password is required")
)

type User struct {
	ID       entity.ID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

func NewUser(name, email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:       entity.NewID(),
		Name:     name,
		Email:    email,
		Password: string(hash),
	}

	err = user.Validate(password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Validate(password string) error {
	if u.ID.String() == "" {
		return ErrIDIsRequired
	}
	if _, err := entity.ParseID(u.ID.String()); err != nil {
		return ErrInvalidId
	}
	if u.Name == "" {
		return ErrNameIsRequired
	}
	if u.Email == "" {
		return ErrEmailIsRequired
	}
	if password == "" {
		return ErrPasswordIsRequired
	}
	return nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}
