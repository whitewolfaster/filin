package models

import (
	"crypto/sha256"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"`
	EncPassword string `json:"-"`
	Token       string `json:"-"`
	Activated   int    `json:"-"`
}

func (u *User) BeforeCreate() error {
	if err := validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 100)),
	); err != nil {
		return err
	}
	if len(u.Password) > 0 {
		u.EncPassword = EncryptString(u.Password)
	}

	return nil
}

func (u *User) ComparePassword(password string) bool {
	sum := sha256.Sum256([]byte(password))
	pswd := fmt.Sprintf("%x", sum)

	if pswd == u.EncPassword {
		return true
	}
	return false
}

func EncryptString(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}
