package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/assert/v2"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

func TestHandler_registerUser(t *testing.T) {
	tests := []struct {
		err        error
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "simple test #1",
			statusCode: 400,
		},
		{
			name:       "simple test #2",
			body:       `{"login": "login", "password": "password"}`,
			statusCode: 200,
		},
		{
			name:       "simple test #3",
			body:       `{"login": "login", "password": "password"}`,
			err:        errors.New("error"),
			statusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				storage: &mockStorage{err: tt.err},
				chOrder: make(chan model.UserOrder),
				token:   jwtauth.New("HS256", []byte("SecretKey"), nil),
			}

			bodyReader := bytes.NewReader([]byte(tt.body))

			r := httptest.NewRequest(http.MethodPost, "/", bodyReader)
			r = r.WithContext(contextWithJWT(r.Context()))
			w := httptest.NewRecorder()

			h.registerUser(w, r)

			result := w.Result()
			result.Body.Close()

			assert.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_loginUser(t *testing.T) {
	tests := []struct {
		err        error
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "simple test #1",
			statusCode: 400,
		},
		{
			name:       "simple test #2",
			body:       `{"login": "login", "password": "password"}`,
			statusCode: 200,
		},
		{
			name:       "simple test #3",
			body:       `{"login": "login", "password": "password"}`,
			err:        errors.New("error"),
			statusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				storage: &mockStorage{err: tt.err},
				chOrder: make(chan model.UserOrder),
				token:   jwtauth.New("HS256", []byte("SecretKey"), nil),
			}

			bodyReader := bytes.NewReader([]byte(tt.body))

			r := httptest.NewRequest(http.MethodPost, "/", bodyReader)
			r = r.WithContext(contextWithJWT(r.Context()))
			w := httptest.NewRecorder()

			h.loginUser(w, r)

			result := w.Result()
			result.Body.Close()

			assert.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}
