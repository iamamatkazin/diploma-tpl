package postgresql

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

func (s *Storage) Run(ctx context.Context) error {
	if err := s.restartPolling(ctx); err != nil {
		return err
	}

	go s.worker(ctx)

	return nil
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
				timerPoll.Reset(time.Second)
				break
			}

			switch code {
			case http.StatusOK:
				var accrual model.Accrual
				if err := json.Unmarshal(data, &accrual); err != nil {
					slog.Error("Ошибка получения данных от системы расчета баллов", slog.Any("error", err))
					return
				}

				if err := s.updateOrder(ctx, accrual, order); err != nil {
					// в идеале функцию изменения данных нужно поместить в отдельный воркер,
					// который будет пытаться до последнего сохранить данные в случае возникновения ошибки
					slog.Error("Ошибка получения данных от системы расчета баллов", slog.Any("error", err))
					timerPoll.Reset(time.Second)
					break
				}

				if accrual.Status == model.Registrered || accrual.Status == model.Processing {
					timerPoll.Reset(time.Second)
				}

			case http.StatusInternalServerError:
				slog.Error("Ошибка получения данных от системы расчета баллов", slog.Any("error", string(data)))
				return

			// если код ответа отличен от 200 и 500, пробуем получить данные снова
			default:
				timerPoll.Reset(time.Second)
			}
		}
	}
}

// Запускаем опрос заказов, которые по какой-то причине еще не прошли систему
// получения баллов.
func (s *Storage) restartPolling(ctx context.Context) error {
	list, err := s.loadUnprocessedOrders(ctx)
	if err != nil {
		return err
	}

	for _, item := range list {
		s.chOrder <- item
	}

	return nil
}
