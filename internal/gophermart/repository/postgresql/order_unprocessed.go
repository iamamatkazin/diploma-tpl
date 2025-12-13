package postgresql

import (
	"context"

	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

func (s *Storage) loadUnprocessedOrders(ctx context.Context) ([]model.UserOrder, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		SELECT login, number FROM orders
		WHERE status != 'PROCESSED' AND status != 'INVALID' 
	`

	rows, err := retryableQuery(ctx, tx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// как быть с емкостью слайса? нужно ли сначала делать запрос на получение количесва записей?
	orders := make([]model.UserOrder, 0)
	for rows.Next() {
		var order model.UserOrder
		err = rows.Scan(&order.Login, &order.Order)
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
