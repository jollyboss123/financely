package postgresstore

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/cespare/xxhash/v2"
	"github.com/google/uuid"
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"github.com/jollyboss123/finance-tracker/pkg/middleware"
	"time"
)

type PostgresStore struct {
	db          *sql.DB
	stopCleanup chan bool

	logger *logger.Logger
}

func New(db *sql.DB, l *logger.Logger) *PostgresStore {
	return NewWithCleanupInterval(db, l, 5*time.Minute)
}

func NewWithCleanupInterval(db *sql.DB, l *logger.Logger, cleanupInterval time.Duration) *PostgresStore {
	p := &PostgresStore{db: db, logger: l}
	if cleanupInterval > 0 {
		go p.startCleanup(cleanupInterval)
	}
	return p
}

func (p *PostgresStore) Delete(token string) (err error) {
	panic("missing context arg")
}

func (p *PostgresStore) Find(token string) (b []byte, found bool, err error) {
	panic("missing context arg")
}

func (p *PostgresStore) Commit(token string, b []byte, expiry time.Time) (err error) {
	panic("missing context arg")
}

func (p *PostgresStore) DeleteCtx(ctx context.Context, token string) error {
	hash, err := sum(token)
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, "delete from sessions where token = $1", hash)
	return err
}

func (p *PostgresStore) FindCtx(ctx context.Context, token string) (b []byte, found bool, err error) {
	hash, err := sum(token)
	if err != nil {
		return nil, false, err
	}

	if err := p.db.QueryRowContext(ctx, `
select data from sessions
where token = $1
and current_timestamp < expiry
order by expiry desc`, hash).Scan(&b); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return b, true, nil
}

func (p *PostgresStore) CommitCtx(ctx context.Context, token string, b []byte, expiry time.Time) error {
	var userID any
	userID, ok := ctx.Value(middleware.KeyID).(uuid.UUID)
	if !ok {
		userID = nil
	}

	hash, err := sum(token)
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, `
insert into sessions (token, user_id, data, expiry)
values ($1, $2, $3, $4)
on conflict (token)
do update set data = excluded.data,
    expiry = excluded.expiry`, hash, userID, b, expiry)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresStore) startCleanup(interval time.Duration) {
	p.stopCleanup = make(chan bool)
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			err := p.deleteExpired()
			if err != nil {
				p.logger.Error().Err(err).Msg("failed delete expired session")
			}
		case <-p.stopCleanup:
			ticker.Stop()
			return
		}
	}
}

func (p *PostgresStore) StopCleanup() {
	if p.stopCleanup != nil {
		p.stopCleanup <- true
	}
}

func (p *PostgresStore) deleteExpired() error {
	_, err := p.db.Exec("delete from sessions where expiry < current_timestamp")
	return err
}

func sum(token string) (string, error) {
	h := xxhash.New()
	_, err := h.Write([]byte(token))
	if err != nil {
		return "", err
	}

	sum := h.Sum(nil)
	str := hex.EncodeToString(sum)
	return str, nil
}
