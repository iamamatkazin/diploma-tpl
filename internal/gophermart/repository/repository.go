package repository

import (
	"context"

	"github.com/iamamatkazin/diploma-tpl/internal/config"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/repository/postgresql"
)

type Storager interface {
	LoadOrder(ctx context.Context, login, order string) (string, error)
	ListOrders(ctx context.Context, login string) ([]model.Order, error)
	GetBalance(ctx context.Context, login string) (model.Balance, error)
	WithdrawBalance(ctx context.Context, login string, withdraw model.Withdraw) (int, error)
	ListWithdrawals(ctx context.Context, login string) ([]model.Withdraw, error)
	LoginUser(ctx context.Context, login model.Login) error
	Shutdown()
}

type Storage struct {
	cfg     *config.Config
	storage Storager
}

func New(ctx context.Context, cfg *config.Config) (*Storage, error) {
	dbStor, err := postgresql.New(cfg)
	if err != nil {
		return nil, err
	}

	if err = dbStor.Run(ctx); err != nil {
		return nil, err
	}

	s := &Storage{
		cfg:     cfg,
		storage: dbStor,
	}

	return s, nil
}

func (s *Storage) Shutdown() {
	s.storage.Shutdown()
}

func (s *Storage) LoadOrder(ctx context.Context, login, order string) (string, error) {
	return s.storage.LoadOrder(ctx, login, order)
}

func (s *Storage) ListOrders(ctx context.Context, login string) ([]model.Order, error) {
	return s.storage.ListOrders(ctx, login)
}

func (s *Storage) GetBalance(ctx context.Context, login string) (model.Balance, error) {
	return s.storage.GetBalance(ctx, login)
}

func (s *Storage) WithdrawBalance(ctx context.Context, login string, withdraw model.Withdraw) (int, error) {
	return s.storage.WithdrawBalance(ctx, login, withdraw)
}

func (s *Storage) ListWithdrawals(ctx context.Context, login string) ([]model.Withdraw, error) {
	return s.storage.ListWithdrawals(ctx, login)
}

func (s *Storage) LoginUser(ctx context.Context, login model.Login) error {
	return nil
}
