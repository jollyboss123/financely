package expense

import (
	"database/sql"
	s "github.com/shopspring/decimal"
	"time"
)

type Schema struct {
	ID              int          `db:"id"`
	Title           string       `db:"title"`
	Amount          s.Decimal    `db:"amount"`
	TransactionDate time.Time    `db:"transaction_date"`
	CreatedAt       time.Time    `db:"created_at"`
	UpdatedAt       time.Time    `db:"updated_at"`
	DeletedAt       sql.NullTime `db:"deleted_at"`
}
