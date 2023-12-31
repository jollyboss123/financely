package database

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/jollyboss123/finance-tracker/config"
	"github.com/jollyboss123/finance-tracker/pkg/logger"
	"log"
	"math/rand"
	"time"
)

func New(l *logger.Logger, cfg config.Database) *sqlx.DB {
	var dsn string
	dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name)
	cfg.Driver = "pgx"

	db, err := sqlx.Open(cfg.Driver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	alive(l, db.DB)

	return db
}

func alive(l *logger.Logger, db *sql.DB) {
	l.Info().Msg("Connecting to database...")
	base, capacity := time.Second, time.Minute
	backoff := base

	for {
		_, err := db.Exec("SELECT true")
		if err == nil {
			l.Info().Msg("Database connected")
			return
		}

		l.Error().Err(err).Msg("Database connection failed")

		jitter := time.Duration(rand.Int63n(int64(backoff * 3 / 2)))
		sleep := base + jitter
		time.Sleep(sleep)

		if backoff < capacity {
			backoff <<= 1
		}
	}
}
