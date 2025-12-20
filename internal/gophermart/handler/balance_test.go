package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

func TestHandler_getBalance(t *testing.T) {
	tests := []struct {
		err        error
		name       string
		statusCode int
	}{
		{
			name:       "simple test #1",
			statusCode: 200,
		},
		{
			name:       "simple test #2",
			err:        errors.New("error"),
			statusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				storage: &mockStorage{err: tt.err},
				chOrder: make(chan model.UserOrder),
			}

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r = r.WithContext(contextWithJWT(r.Context()))
			w := httptest.NewRecorder()

			h.getBalance(w, r)

			result := w.Result()
			result.Body.Close()

			assert.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_withdrawBalance(t *testing.T) {
	tests := []struct {
		err        error
		name       string
		body       string
		code       int
		statusCode int
	}{
		{
			name:       "simple test #1",
			statusCode: 400,
		},
		{
			name:       "simple test #2",
			body:       `{"order": "12345678903", "sum": 10}`,
			code:       200,
			statusCode: 200,
		},
		{
			name:       "simple test #3",
			body:       `{"order": "12345678903", "sum": 10}`,
			code:       402,
			statusCode: 402,
		},
		{
			name:       "simple test #4",
			body:       `{"order": "12345678903", "sum": 10}`,
			err:        errors.New("error"),
			statusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				storage: &mockStorage{code: tt.code, err: tt.err},
				chOrder: make(chan model.UserOrder),
			}

			bodyReader := bytes.NewReader([]byte(tt.body))

			r := httptest.NewRequest(http.MethodPost, "/", bodyReader)
			r = r.WithContext(contextWithJWT(r.Context()))
			w := httptest.NewRecorder()

			h.withdrawBalance(w, r)

			result := w.Result()
			result.Body.Close()

			assert.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_listWithdrawals(t *testing.T) {
	tests := []struct {
		err        error
		name       string
		list       []model.Withdraw
		statusCode int
	}{
		{
			name:       "simple test #1",
			statusCode: 204,
		},
		{
			name:       "simple test #2",
			list:       []model.Withdraw{{}},
			statusCode: 200,
		},
		{
			name:       "simple test #3",
			list:       []model.Withdraw{{}},
			err:        errors.New("error"),
			statusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				storage: &mockStorage{listWithdraw: tt.list, err: tt.err},
				chOrder: make(chan model.UserOrder),
			}

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r = r.WithContext(contextWithJWT(r.Context()))
			w := httptest.NewRecorder()

			h.listWithdrawals(w, r)

			result := w.Result()
			result.Body.Close()

			assert.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}
