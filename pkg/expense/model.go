package expense

import (
	"github.com/google/uuid"
	"time"
)

type Schema struct {
	ID              uuid.UUID `db:"id"`
	Title           string    `db:"title"`
	Amount          uint64    `db:"amount"`
	CurrencyID      uuid.UUID `db:"currency_id"`
	TransactionDate time.Time `db:"transaction_date"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
