package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
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

	select {
	case s.chOrder <- model.UserOrder{Login: login, Order: order}:
	default:
		slog.Info("Занят канал отправки заказа в систему расчета начислений")
	}
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
	var accrual sql.NullFloat64

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
		fmt.Println("&&&&& retryableQuery", err)
		return nil, err
	}
	defer rows.Close()

	// как быть с емкостью слайса? нужно ли сначала делать запрос на получение количесва записей?
	orders := make([]model.Order, 0)
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.Number, &order.Status, &accrual, &order.UploadedAt)
		if err != nil {
			fmt.Println("&&&&& rows.Scan", err)
			return nil, err
		}

		if accrual.Valid {
			order.Accrual = accrual.Float64 / 100
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println("&&&&& rows.Err", err)
		return nil, err
	}

	tx.Commit()

	return orders, nil
}
