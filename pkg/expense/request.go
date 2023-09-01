package expense

import (
	s "github.com/shopspring/decimal"
	"time"
)

type CreateRequest struct {
	Title           string    `json:"title" validate:"required"`
	Amount          s.Decimal `json:"amount" validate:"required"`
	TransactionDate time.Time `json:"transaction_date" validate:"required"`
}

type UpdateRequest struct {
	ID              int       `json:"-"`
	Title           string    `json:"title" validate:"required"`
	Amount          s.Decimal `json:"amount" validate:"required"`
	TransactionDate time.Time `json:"transaction_date" validate:"required"`
}
