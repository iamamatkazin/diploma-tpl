package postgresql

import (
	"context"
	"database/sql"

	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

func (s *Storage) UpdateOrder(ctx context.Context, accrual model.Accrual, order model.UserOrder) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if accrual.Status == model.Processed {
		if err := updateUserCurrent(ctx, tx, accrual, order); err != nil {
			return err
		}
	}

	if err := updateOrderStatus(ctx, tx, accrual); err != nil {
		return err
	}

	return tx.Commit()
}

func updateOrderStatus(ctx context.Context, tx *sql.Tx, accrual model.Accrual) error {
	query := `
		UPDATE orders SET 
			status = $2,
			accrual =$3
    	WHERE number = $1
	`

	return retryableExec(ctx, tx, query, accrual.Order, accrual.Status, accrual.Accrual)
}

func updateUserCurrent(ctx context.Context, tx *sql.Tx, accrual model.Accrual, order model.UserOrder) error {
	query := `
		UPDATE users SET 
			current = users.current + $2
    	WHERE login = $1
	`

	return retryableExec(ctx, tx, query, order.Login, accrual.Accrual)
}
