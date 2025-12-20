package handler

import (
	"context"

	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

type mockStorage struct {
	err          error
	login        string
	listOrder    []model.Order
	listWithdraw []model.Withdraw
	code         int
}

func (s *mockStorage) Shutdown() {
}

func (s *mockStorage) LoadOrder(ctx context.Context, login, order string) (string, error) {
	return s.login, s.err
}

func (s *mockStorage) ListOrders(ctx context.Context, login string) ([]model.Order, error) {
	return s.listOrder, s.err
}

func (s *mockStorage) GetBalance(ctx context.Context, login string) (model.Balance, error) {
	return model.Balance{}, s.err
}

func (s *mockStorage) WithdrawBalance(ctx context.Context, login string, withdraw model.Withdraw) (int, error) {
	return s.code, s.err
}

func (s *mockStorage) ListWithdrawals(ctx context.Context, login string) ([]model.Withdraw, error) {
	return s.listWithdraw, s.err
}

func (s *mockStorage) UpdateOrder(ctx context.Context, accrual model.Accrual, order model.UserOrder) error {
	return s.err
}

func (s *mockStorage) LoadUnprocessedOrders(ctx context.Context) ([]model.UserOrder, error) {
	return nil, s.err
}

func (s *mockStorage) LoginUser(ctx context.Context, login model.Login) (*model.Login, error) {
	return nil, s.err
}

func (s *mockStorage) RegisterUser(ctx context.Context, login model.Login) (*model.Login, error) {
	return nil, s.err
}
