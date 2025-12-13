package accrual

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/iamamatkazin/diploma-tpl/internal/config"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

const timeout = time.Second * 5

type Accrual struct {
	client *http.Client
	cfg    *config.Config
}

func New(cfg *config.Config) *Accrual {
	return &Accrual{
		cfg: cfg,
		client: &http.Client{
			Timeout:   timeout,
			Transport: &http.Transport{},
		},
	}
}

func (a *Accrual) Get(ctx context.Context, order model.UserOrder) (data []byte, code int, err error) {
	url := fmt.Sprintf("%s/api/orders/%s", a.cfg.AccAddress, order.Order)

	data, code, err = a.get(ctx, url, "application/json")
	if err != nil {
		return nil, 0, err
	}

	fmt.Println("##########", url, code, string(data))
	return data, code, nil
}

func (a *Accrual) get(ctx context.Context, url, contentType string) (data []byte, code int, err error) {
	request, err := newRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, 0, err
	}

	request.Header.Set("Content-Type", contentType)

	// отправляем запрос и получаем ответ
	response, err := a.client.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()

	data, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, 0, err
	}

	return data, response.StatusCode, nil
}

func newRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	timerRetriable := time.NewTimer(0)
	count := 0

	for {
		select {
		case <-ctx.Done():
		case <-timerRetriable.C:
			request, err := http.NewRequestWithContext(ctx, method, url, body)
			if err != nil {
				if count > 3 {
					return nil, err
				}

				timerRetriable.Reset(time.Duration(2*count+1) * time.Second)
				count++
			}

			return request, nil
		}
	}
}
