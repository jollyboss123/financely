package expense

import (
	"github.com/google/uuid"
	"time"
)

type CreateRequest struct {
	Title           string    `json:"title" validate:"required"`
	Amount          uint64    `json:"amount" validate:"required,gt=0"`
	CurrencyCode    string    `json:"currency_code" validate:"required"`
	CurrencyID      uuid.UUID `json:"currency_id" validate:"required"`
	TransactionDate time.Time `json:"transaction_date" validate:"required"`
}

type UpdateRequest struct {
	ID              uuid.UUID `json:"-"`
	Title           string    `json:"title"`
	Amount          uint64    `json:"amount" validate:"gt=0"`
	CurrencyCode    string    `json:"currency_code"`
	CurrencyID      uuid.UUID `json:"currency_id"`
	TransactionDate time.Time `json:"transaction_date" validate:"required"`
}
