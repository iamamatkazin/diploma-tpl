package accrual

import (
	"sync"

	"github.com/iamamatkazin/diploma-tpl/internal/agent"
	"github.com/iamamatkazin/diploma-tpl/internal/config"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/repository"
)

type Accrual struct {
	agent   *agent.Agent
	chOrder chan model.UserOrder
	storage *repository.Storage
	sync.WaitGroup
}

func New(cfg *config.Config, chOrder chan model.UserOrder, storage *repository.Storage) *Accrual {
	return &Accrual{
		chOrder: chOrder,
		agent:   agent.New(cfg),
		storage: storage,
	}
}
