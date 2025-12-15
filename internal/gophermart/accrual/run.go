package accrual

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

func (a *Accrual) Run(ctx context.Context) error {
	a.Add(1)
	go func() {
		defer a.Done()
		a.worker(ctx)
	}()

	if err := a.restartPolling(ctx); err != nil {
		return err
	}

	return nil
}

// Создаем воркер, который получает заказ и пытается для него получить расчет баллов.
// Если баллы удалось рассчитать или расчет не возможен,
// то заказ изменяет свои поля в базе данных.
func (a *Accrual) worker(ctx context.Context) {
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return

		case userOrder := <-a.chOrder:
			// что бы не реализовывать буфферезированный канал не понятно какого размера,
			// сразу помещаем заказ в горутину, которая после получения расчета завершится сама.
			wg.Add(1)
			go func(order model.UserOrder) {
				defer wg.Done()
				a.analyzeResponse(ctx, order)
			}(userOrder)
		}
	}
}

func (a *Accrual) analyzeResponse(ctx context.Context, order model.UserOrder) {
	timerPoll := time.NewTimer(0)
	defer timerPoll.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-timerPoll.C:
			data, code, err := a.agent.Get(ctx, order)
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

				if accrual.Accrual != nil {
					val := *accrual.Accrual * 100
					accrual.Accrual = &val
				}

				if err := a.storage.UpdateOrder(ctx, accrual, order); err != nil {
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
func (a *Accrual) restartPolling(ctx context.Context) error {
	list, err := a.storage.LoadUnprocessedOrders(ctx)
	if err != nil {
		return err
	}

	for _, item := range list {
		a.chOrder <- item
	}

	return nil
}
