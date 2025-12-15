package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/iamamatkazin/diploma-tpl/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	cfg *config.Config
	db  *sql.DB
}

func New(cfg *config.Config) (*Storage, error) {
	if cfg.Database == "" {
		return &Storage{cfg: cfg}, fmt.Errorf("не возможно соединиться с базой данных")
	}

	db, err := sql.Open("pgx", cfg.Database)
	if err != nil {
		return nil, err
	}

	loadMigrations(db)

	return &Storage{
		cfg: cfg,
		db:  db,
	}, nil
}

func loadMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error(err.Error())
	}
}

func isRetryablePgError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	switch pgErr.Code {
	case pgerrcode.ConnectionException,
		pgerrcode.ConnectionDoesNotExist,
		pgerrcode.ConnectionFailure,
		pgerrcode.SQLClientUnableToEstablishSQLConnection,
		pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection,
		pgerrcode.TransactionResolutionUnknown,
		pgerrcode.ProtocolViolation:
		return true
	default:
		return false
	}
}

func retryableExec(ctx context.Context, tx *sql.Tx, query string, args ...any) error {
	timerRetryable := time.NewTimer(0)
	count := 0

	for {
		select {
		case <-ctx.Done():
		case <-timerRetryable.C:
			_, err := tx.ExecContext(ctx, query, args...)
			if err != nil {
				if !isRetryablePgError(err) || count > 3 {
					return err
				}

				timerRetryable.Reset(time.Duration(2*count+1) * time.Second)
				count++
			}

			return nil
		}
	}
}

func retryableQuery(ctx context.Context, tx *sql.Tx, query string, args ...any) (*sql.Rows, error) {
	var (
		rows *sql.Rows
		err  error
	)

	timerRetryable := time.NewTimer(0)
	count := 0

	for {
		select {
		case <-ctx.Done():
		case <-timerRetryable.C:
			if len(args) == 0 {
				rows, err = tx.QueryContext(ctx, query)
			} else {
				rows, err = tx.QueryContext(ctx, query, args...)
			}
			if err != nil {
				if !isRetryablePgError(err) || count > 3 {
					return nil, err
				}

				timerRetryable.Reset(time.Duration(2*count+1) * time.Second)
				count++
			}

			return rows, nil
		}
	}
}

func (s *Storage) Shutdown() {
	if s.db != nil {
		s.db.Close()
	}
}
