package currency

import (
	"github.com/google/uuid"
)

type CreateRequest struct {
	Code        string `json:"code" validate:"required"`
	NumericCode string `json:"numeric_code" validate:"required"`
	Fraction    int8   `json:"fraction" validate:"required,gt=0"`
	Grapheme    string `json:"grapheme" validate:"required"`
	Template    string `json:"template" validate:"required"`
	Decimal     string `json:"decimal"`
	Thousand    string `json:"thousand"`
}

type UpdateRequest struct {
	ID          uuid.UUID `json:"-"`
	Code        string    `json:"code"`
	NumericCode string    `json:"numeric_code"`
	Fraction    int8      `json:"fraction" validate:"gt=0"`
	Grapheme    string    `json:"grapheme"`
	Template    string    `json:"template"`
	Decimal     string    `json:"decimal"`
	Thousand    string    `json:"thousand"`
}
