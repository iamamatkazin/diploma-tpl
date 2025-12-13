package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

func (s *Storage) LoginUser(ctx context.Context, login model.Login) (*model.Login, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := loginUser(ctx, tx, login)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return user, nil
}

func loginUser(ctx context.Context, tx *sql.Tx, login model.Login) (*model.Login, error) {
	query := `
		SELECT login, password FROM users
		WHERE login = $1 
	`

	var user model.Login
	if err := tx.QueryRowContext(ctx, query, login.Login).Scan(&user.Login, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (s *Storage) RegisterUser(ctx context.Context, login model.Login) (*model.Login, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := loginUser(ctx, tx, login)
	if err != nil {
		return nil, err
	}

	if user == nil {
		if err = registerUser(ctx, tx, login); err != nil {
			return nil, err
		}
	}

	tx.Commit()

	return user, nil
}

func registerUser(ctx context.Context, tx *sql.Tx, login model.Login) error {
	query := `
		INSERT INTO users (login, password, current, withdrawn)
		VALUES ($1, $2, 0, 0)
	`

	if err := retryableExec(ctx, tx, query, login.Login, login.Password); err != nil {
		return err
	}

	return nil
}
