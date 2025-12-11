package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

func (s *Storage) LoadOrder(ctx context.Context, login, order string) (currentLogin string, err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	currentLogin, err = getUserByOrder(ctx, tx, order)
	if currentLogin == "" {
		if err = loadOrder(ctx, tx, login, order); err != nil {
			return "", err
		}
	}

	tx.Commit()

	return currentLogin, nil
}

func getUserByOrder(ctx context.Context, tx *sql.Tx, order string) (string, error) {
	query := `
		SELECT login FROM orders
		WHERE number = $1 
	`

	var login string
	err := tx.QueryRowContext(ctx, query, order).Scan(&login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return login, nil
}

func loadOrder(ctx context.Context, tx *sql.Tx, login, order string) error {
	query := `
		INSERT INTO orders (number, login, status, uploaded_at)
		VALUES ($1, $2, $3, $4)
	`
	return retryableExec(ctx, tx, query, order, login, model.New, time.Now().Local())
}

func (s *Storage) ListOrders(ctx context.Context, login string) ([]model.Order, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		SELECT number, status, accrual, uploaded_at FROM orders
		WHERE login = $1 
	`

	rows, err := retryableQuery(ctx, tx, query, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// как быть с емкостью слайса? нужно ли сначала делать запрос на получение количесва записей?
	orders := make([]model.Order, 0)
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return orders, nil
}
