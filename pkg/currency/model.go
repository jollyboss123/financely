package currency

import (
	"github.com/google/uuid"
	"time"
)

type Schema struct {
	ID          uuid.UUID `db:"id"`
	Code        string    `db:"code"`
	NumericCode string    `db:"numeric_code"`
	Fraction    int8      `db:"fraction"`
	Grapheme    string    `db:"grapheme"`
	Template    string    `db:"template"`
	Decimal     string    `db:"decimal"`
	Thousand    string    `db:"thousand"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
