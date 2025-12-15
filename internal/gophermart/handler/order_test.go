package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/assert/v2"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Вспомогательная функция: создаёт контекст с токеном
func contextWithJWT(ctx context.Context) context.Context {
	claims := map[string]interface{}{
		"login": "login",
	}

	tokenAuth := jwtauth.New("HS256", []byte("SecretKey"), nil)
	_, tokenString, _ := tokenAuth.Encode(claims)
	parsedToken, _ := jwt.ParseString(tokenString, jwt.WithVerify(false))

	ctx = context.WithValue(ctx, jwtauth.TokenCtxKey, parsedToken)

	chiCtx := chi.NewRouteContext()
	ctx = context.WithValue(ctx, chi.RouteCtxKey, chiCtx)

	return ctx
}

func TestHandler_loadOrder(t *testing.T) {
	tests := []struct {
		err        error
		name       string
		body       string
		login      string
		statusCode int
	}{
		{
			name:       "simple test #1",
			body:       "12345678903",
			statusCode: 202,
		},
		{
			name:       "simple test #2",
			body:       "12345678903",
			login:      "login",
			statusCode: 200,
		},
		{
			name:       "simple test #3",
			body:       "12345678903",
			login:      "login111",
			statusCode: 409,
		},
		{
			name:       "simple test #4",
			body:       "12345678903",
			login:      "login",
			err:        errors.New("error"),
			statusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				storage: &mockStorage{login: tt.login, err: tt.err},
				chOrder: make(chan model.UserOrder),
			}

			r := httptest.NewRequest(http.MethodPost, "/", nil)
			r = r.WithContext(contextWithJWT(r.Context()))
			w := httptest.NewRecorder()

			h.loadOrder(w, r)

			result := w.Result()
			result.Body.Close()

			assert.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}

func TestHandler_listOrders(t *testing.T) {
	tests := []struct {
		err        error
		name       string
		list       []model.Order
		statusCode int
	}{
		{
			name:       "simple test #1",
			statusCode: 204,
		},
		{
			name:       "simple test #2",
			list:       []model.Order{{}},
			statusCode: 200,
		},
		{
			name:       "simple test #3",
			err:        errors.New("error"),
			statusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				storage: &mockStorage{listOrder: tt.list, err: tt.err},
				chOrder: make(chan model.UserOrder),
			}

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r = r.WithContext(contextWithJWT(r.Context()))
			w := httptest.NewRecorder()

			h.listOrders(w, r)

			result := w.Result()
			result.Body.Close()

			assert.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}
