package expense

import s "github.com/shopspring/decimal"

type CreateRequest struct {
	Title  string    `json:"title" validate:"required"`
	Amount s.Decimal `json:"amount" validate:"required"`
}

type UpdateRequest struct {
	ID     int       `json:"-"`
	Title  string    `json:"title" validate:"required"`
	Amount s.Decimal `json:"amount" validate:"required"`
}
