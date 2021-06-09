package models

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

var (
	ErrInvalidBookPriceOrYear = errors.New("invalid book price or book year value")
)

type Book struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Author      string `json:"author"`
	Year        int    `json:"year"`
	Genre       string `json:"genre"`
	PubHouse    string `json:"pub_house"`
	CoverPath   string `json:"cover_path"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type Genre struct {
	Name string `json:"name"`
}

func (b *Book) BeforeCreate() error {
	if err := validation.ValidateStruct(
		b,
		validation.Field(&b.Name, validation.Required),
		validation.Field(&b.Author, validation.Required),
		validation.Field(&b.Year, validation.Required),
		validation.Field(&b.Genre, validation.Required),
		validation.Field(&b.PubHouse, validation.Required),
		validation.Field(&b.Description, validation.Required),
		validation.Field(&b.Price, validation.Required),
	); err != nil {
		return err
	}
	if b.Price <= 0 || b.Year <= 0 {
		return ErrInvalidBookPriceOrYear
	}
	return nil
}

func (b *Book) BeforeUpdate() error {
	if err := validation.ValidateStruct(
		b,
		validation.Field(&b.ID, validation.Required),
		validation.Field(&b.Name, validation.Required),
		validation.Field(&b.Author, validation.Required),
		validation.Field(&b.Year, validation.Required),
		validation.Field(&b.Genre, validation.Required),
		validation.Field(&b.PubHouse, validation.Required),
		validation.Field(&b.Description, validation.Required),
		validation.Field(&b.Price, validation.Required),
	); err != nil {
		return err
	}
	if b.Price <= 0 || b.Year <= 0 {
		return ErrInvalidBookPriceOrYear
	}
	return nil
}
