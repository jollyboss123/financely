package authentication

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/jollyboss123/finance-tracker/pkg/server/message"
	"github.com/jollyboss123/scs/v2"
	"time"
)

type User interface {
	Register(ctx context.Context, firstName, lastName, email, hashedPassword string) error
	Login(ctx context.Context, req *LoginRequest) (*Schema, bool, error)
	Logout(ctx context.Context, userID uuid.UUID) (bool, error)
	Csrf(ctx context.Context) (string, error)
}

type userRepository struct {
	db      *sqlx.DB
	session *scs.SessionManager
}

const (
	InsertIntoUsers = `insert into users
(first_name, last_name, email, password)
values ($1, $2, $3, $4) returning id`
	SelectUserByEmail = "select * from users where lower(email) = lower($1)"
	CheckUserSession  = `select case 
when exists(select * 
	from sessions
	where sessions.user_id = $1)
then true
else false
end`
	DeleteUserSession = `delete from sessions where user_id = $1`
)

var (
	ErrEmailNotAvailable = errors.New("email is not available")
	ErrNotLoggedIn       = errors.New("you are not logged in yet")
)

func New(db *sqlx.DB, manager *scs.SessionManager) *userRepository {
	return &userRepository{
		db:      db,
		session: manager,
	}
}

func (ur *userRepository) Register(ctx context.Context, firstName, lastName, email, hashedPassword string) error {
	var userID uuid.UUID

	if err := ur.db.QueryRowContext(ctx, InsertIntoUsers,
		firstName,
		lastName,
		email,
		hashedPassword).Scan(&userID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return ErrEmailNotAvailable
			}
		}

		return err
	}

	return nil
}

func (ur *userRepository) Login(ctx context.Context, req *LoginRequest) (*Schema, bool, error) {
	var u Schema
	err := ur.db.GetContext(ctx, &u, SelectUserByEmail, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, message.ErrBadRequest
		}
		return nil, false, err
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, u.Password)
	if err != nil {
		return nil, false, errors.New("wrong password is provided")
	}

	return &u, match, nil
}

func (ur *userRepository) Logout(ctx context.Context, userID uuid.UUID) (bool, error) {
	var found bool
	if err := ur.db.QueryRowContext(ctx, CheckUserSession, userID).Scan(&found); err != nil {
		return false, err
	}

	if !found {
		return false, ErrNotLoggedIn
	}

	_, err := ur.db.ExecContext(ctx, DeleteUserSession, userID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (ur *userRepository) Csrf(ctx context.Context) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	err = ur.session.CtxStore.CommitCtx(ctx, token, []byte("csrf_token"), time.Now().Add(ur.session.Lifetime))
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
