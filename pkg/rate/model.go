package rate

import (
	"github.com/google/uuid"
	"time"
)

type Schema struct {
	From      uuid.UUID `db:"from_currency_id"`
	To        uuid.UUID `db:"to_currency_id"`
	Rate      float64   `db:"rate"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
