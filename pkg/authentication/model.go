package authentication

import (
	"github.com/google/uuid"
	"time"
)

type Schema struct {
	ID         uuid.UUID  `db:"id"`
	FirstName  string     `db:"first_name"`
	LastName   string     `db:"last_name"`
	Email      string     `db:"email"`
	Password   string     `db:"password"`
	VerifiedAt *time.Time `db:"verified_at"`
}
