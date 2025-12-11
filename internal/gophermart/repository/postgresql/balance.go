package postgresql

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

func (s *Storage) GetBalance(ctx context.Context, login string) (model.Balance, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Balance{}, err
	}
	defer tx.Rollback()

	balance, err := getBalance(ctx, tx, login)
	if err != nil {
		return model.Balance{}, err
	}

	tx.Commit()

	return balance, nil
}

func getBalance(ctx context.Context, tx *sql.Tx, login string) (model.Balance, error) {
	query := `
		SELECT current, withdrawn FROM users
		WHERE login = $1 
	`
	var balance model.Balance
	err := tx.QueryRowContext(ctx, query, login).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return model.Balance{}, err
	}

	return balance, nil
}

func (s *Storage) WithdrawBalance(ctx context.Context, login string, withdraw model.Withdraw) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	balance, err := getBalance(ctx, tx, login)
	if err != nil {
		return 0, err
	}

	if balance.Current < withdraw.Sum {
		return http.StatusPaymentRequired, nil
	}

	if err = updateBalance(ctx, tx, login, withdraw.Sum); err != nil {
		return 0, err
	}

	if err := updateOrderSum(ctx, tx, withdraw); err != nil {
		return 0, err
	}

	tx.Commit()

	return http.StatusOK, nil
}

func updateOrderSum(ctx context.Context, tx *sql.Tx, withdraw model.Withdraw) error {
	query := `
		UPDATE orders SET 
			sum = $2,
			processed_at =$3
    	WHERE number = $1
	`

	return retryableExec(ctx, tx, query, withdraw.Order, withdraw.Sum, time.Now().Local())
}

func updateBalance(ctx context.Context, tx *sql.Tx, login string, sum int) error {
	query := `
		UPDATE users SET 
			current = users.current - $2,
			withdrawn = users.withdrawn + $3
    	WHERE login = $1
	`

	return retryableExec(ctx, tx, query, login, sum, sum)
}

func (s *Storage) ListWithdrawals(ctx context.Context, login string) ([]model.Withdraw, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		SELECT number, sum, processed_at FROM orders
		WHERE login = $1 AND sum > 0
	`

	rows, err := retryableQuery(ctx, tx, query, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// как быть с емкостью слайса? нужно ли сначала делать запрос на получение количесва записей?
	orders := make([]model.Withdraw, 0)
	for rows.Next() {
		var order model.Withdraw
		err = rows.Scan(&order.Order, &order.Sum, &order.ProcessedAt)
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
