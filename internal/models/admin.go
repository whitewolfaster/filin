package models

import (
	"crypto/sha256"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Admin struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	Password    string `json:"password,omitempty"`
	EncPassword string `json:"-"`
}

func (a *Admin) BeforeCreate() error {
	if err := validation.ValidateStruct(
		a,
		validation.Field(&a.Login, validation.Required, validation.Length(6, 30)),
		validation.Field(&a.Password, validation.Required, validation.Length(8, 30)),
	); err != nil {
		return err
	}
	if len(a.Password) > 0 {
		a.EncPassword = EncryptString(a.Password)
	}
	return nil
}

func (a *Admin) ComparePassword(password string) bool {
	sum := sha256.Sum256([]byte(password))
	pswd := fmt.Sprintf("%x", sum)

	if pswd == a.EncPassword {
		return true
	}
	return false
}
