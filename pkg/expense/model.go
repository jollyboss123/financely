package expense

import (
	"github.com/google/uuid"
	"time"
)

type Schema struct {
	ID               uuid.UUID `db:"id"`
	Title            string    `db:"title"`
	Amount           int64     `db:"amount_ud"`
	CurrencyID       uuid.UUID `db:"currency_id_ud"`
	CurrencyCode     string    `db:"currency_code_ud"`
	BaseAmount       int64     `db:"amount_base"`
	BaseCurrencyID   uuid.UUID `db:"currency_id_base"`
	BaseCurrencyCode string    `db:"currency_code_base"`
	TransactionDate  time.Time `db:"transaction_date"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
