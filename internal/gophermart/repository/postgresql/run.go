package postgresql

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

func (s *Storage) Run(ctx context.Context) {
	// необходимо реализовать восстановление опроса внешнего сервиса при перезапуске системы

	go s.worker(ctx)
}

// Создаем воркер, который получает заказ и пытается для него получить расчет баллов.
// Если баллы удалось рассчитать или расчет не возможен,
// то заказ изменяет свои поля в базе данных.
func (s *Storage) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case userOrder := <-s.chOrder:
			// что бы не реализовывать буфферезированный канал не понятно какого размера,
			// сразу помещаем заказ в горутину, которая после получения расчета завершится сама.
			go func(order model.UserOrder) {
				s.analyzeResponse(ctx, order)
			}(userOrder)
		}
	}
}

func (s *Storage) analyzeResponse(ctx context.Context, order model.UserOrder) {
	timerPoll := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			return

		case <-timerPoll.C:
			data, code, err := s.agent.Get(ctx, order)
			if err != nil {
				slog.Error("Ошибка получения данных от системы расчета баллов", slog.Any("error", err))
				return
			}

			switch code {
			case http.StatusOK:
				var accrual model.Accrual
				if err := json.Unmarshal(data, &accrual); err != nil {
					slog.Error("Ошибка получения данных от системы расчета баллов", slog.Any("error", err))
					return
				}

				if err := s.updateOrder(ctx, accrual, order); err != nil {
					slog.Error("Ошибка получения данных от системы расчета баллов", slog.Any("error", err))
					return
				}

				if accrual.Status == model.Registrered || accrual.Status == model.Processing {
					timerPoll.Reset(time.Second)
				}

				return

			case http.StatusInternalServerError:
				slog.Error("Ошибка получения данных от системы расчета баллов", slog.Any("error", string(data)))
				return

			default:
				timerPoll.Reset(time.Second)
			}
		}
	}
}
