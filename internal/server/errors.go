package server

import (
	"errors"
)

var (
	ErrIncorrectEmailOrPassword = errors.New("incorrect email or password")
	ErrNullContext              = errors.New("context is nil")
	ErrAuthenticated            = errors.New("already authenticated")
	ErrEmailNotSent             = errors.New("email not sent")
	ErrNullProducts             = errors.New("products is nil")
	ErrWrongPassword            = errors.New("wrong old password")
	ErrInvalidBookID            = errors.New("invalid book ID")
	ErrWrongPrimaryKey          = errors.New("wrong primary key")
	ErrExistInCart              = errors.New("product exist in cart")
	ErrInvalidToken             = errors.New("invalid token")
	ErrActivated                = errors.New("already activated")
	ErrNotActivated             = errors.New("not activated")
)
