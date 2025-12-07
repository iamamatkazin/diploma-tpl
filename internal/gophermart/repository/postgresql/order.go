package postgresql

import (
	"context"
	"database/sql"

	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

func (s *Storage) updateOrder(ctx context.Context, accrual model.Accrual, order model.UserOrder) error {
	if s.db == nil {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if accrual.Status == model.Processed {
		if err := s.updateUserCurrent(ctx, tx, accrual, order); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := s.updateOrderStatus(ctx, tx, accrual); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Storage) updateOrderStatus(ctx context.Context, tx *sql.Tx, accrual model.Accrual) error {
	query := `
		UPDATE orders SET 
			status = $2
			accrual =$3
    	WHERE order = $1
	`

	return retryableExec(ctx, tx, query, accrual.Order, accrual.Status, accrual.Accrual)
}

func (s *Storage) updateUserCurrent(ctx context.Context, tx *sql.Tx, accrual model.Accrual, order model.UserOrder) error {
	query := `
		UPDATE users SET 
			current = users.current + $2
    	WHERE login = $1
	`

	return retryableExec(ctx, tx, query, order.Login, accrual.Accrual)
}
