package repository

import (
	"context"
	"net/http"

	"github.com/iamamatkazin/diploma-tpl/internal/config"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/repository/postgresql"
	"github.com/iamamatkazin/diploma-tpl/internal/pkg/custerror"
	"golang.org/x/crypto/bcrypt"
)

type Storager interface {
	LoadOrder(ctx context.Context, login, order string) (string, error)
	ListOrders(ctx context.Context, login string) ([]model.Order, error)
	GetBalance(ctx context.Context, login string) (model.Balance, error)
	WithdrawBalance(ctx context.Context, login string, withdraw model.Withdraw) (int, error)
	ListWithdrawals(ctx context.Context, login string) ([]model.Withdraw, error)
	LoginUser(ctx context.Context, login model.Login) (*model.Login, error)
	RegisterUser(ctx context.Context, login model.Login) (*model.Login, error)
	UpdateOrder(ctx context.Context, accrual model.Accrual, order model.UserOrder) error
	LoadUnprocessedOrders(ctx context.Context) ([]model.UserOrder, error)
	Shutdown()
}

type Storage struct {
	cfg  *config.Config
	stor Storager
}

func New(ctx context.Context, cfg *config.Config, chOrder chan model.UserOrder) (*Storage, error) {
	dbStor, err := postgresql.New(cfg)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		cfg:  cfg,
		stor: dbStor,
	}

	return s, nil
}

func (s *Storage) Shutdown() {
	s.stor.Shutdown()
}

func (s *Storage) LoadOrder(ctx context.Context, login, order string) (string, error) {
	return s.stor.LoadOrder(ctx, login, order)
}

func (s *Storage) ListOrders(ctx context.Context, login string) ([]model.Order, error) {
	return s.stor.ListOrders(ctx, login)
}

func (s *Storage) GetBalance(ctx context.Context, login string) (model.Balance, error) {
	return s.stor.GetBalance(ctx, login)
}

func (s *Storage) WithdrawBalance(ctx context.Context, login string, withdraw model.Withdraw) (int, error) {
	return s.stor.WithdrawBalance(ctx, login, withdraw)
}

func (s *Storage) ListWithdrawals(ctx context.Context, login string) ([]model.Withdraw, error) {
	return s.stor.ListWithdrawals(ctx, login)
}

func (s *Storage) UpdateOrder(ctx context.Context, accrual model.Accrual, order model.UserOrder) error {
	return s.stor.UpdateOrder(ctx, accrual, order)
}

func (s *Storage) LoadUnprocessedOrders(ctx context.Context) ([]model.UserOrder, error) {
	return s.stor.LoadUnprocessedOrders(ctx)
}

func (s *Storage) LoginUser(ctx context.Context, login model.Login) (*model.Login, error) {
	user, err := s.stor.LoginUser(ctx, login)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, custerror.New(http.StatusUnauthorized, "неверная пара логин/пароль")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return nil, custerror.New(http.StatusUnauthorized, "неверная пара логин/пароль")
	}

	return nil, nil
}

func (s *Storage) RegisterUser(ctx context.Context, login model.Login) (*model.Login, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	login.Password = string(hash)

	user, err := s.stor.RegisterUser(ctx, login)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, custerror.New(http.StatusConflict, "логин уже занят")
	}

	return nil, nil
}
