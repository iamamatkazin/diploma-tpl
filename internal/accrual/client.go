package accrual

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const timeout = time.Second * 5

type Client struct {
	*http.Client
}

func New() *Client {
	return &Client{
		Client: &http.Client{
			Timeout:   timeout,
			Transport: &http.Transport{},
		},
	}
}

func (c *Client) Post(ctx context.Context, url, contentType string, data any) (err error) {
	var (
		request *http.Request
		body    []byte
	)

	body, err = json.Marshal(data)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	request, err = newRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	// в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	request.Header.Set("Content-Type", contentType)

	// отправляем запрос и получаем ответ
	response, err := c.Client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка обработки запроса с кодом: %d", response.StatusCode)
	}

	return nil
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
