package repository

import (
	"context"

	"github.com/iamamatkazin/diploma-tpl/internal/config"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/repository/postgresql"
)

type Storager interface {
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

	s := &Storage{
		cfg:     cfg,
		storage: dbStor,
	}

	return s, nil
}

func (s *Storage) Shutdown() {
	s.Shutdown()
}
