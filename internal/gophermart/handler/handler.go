package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"github.com/iamamatkazin/diploma-tpl/internal/config"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/accrual"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/repository"
	"github.com/iamamatkazin/diploma-tpl/internal/pkg/custerror"
)

type Handler struct {
	token   *jwtauth.JWTAuth
	storage repository.Storager
	Router  *chi.Mux
	cfg     *config.Config
	chOrder chan model.UserOrder
	accr    *accrual.Accrual
}

func New(ctx context.Context, cfg *config.Config) (*Handler, error) {
	chOrder := make(chan model.UserOrder)

	storage, err := repository.New(ctx, cfg, chOrder)
	if err != nil {
		return nil, err
	}

	accr := accrual.New(cfg, chOrder, storage)
	if err = accr.Run(ctx); err != nil {
		return nil, err
	}

	h := &Handler{
		storage: storage,
		cfg:     cfg,
		token:   jwtauth.New(cfg.Algorithm, []byte(cfg.SecretKey), nil),
		chOrder: chOrder,
		accr:    accr,
	}

	h.Router = chi.NewRouter()
	h.listRoute()

	return h, nil
}

func (h *Handler) listRoute() {
	// h.Router.Use(middleware.AllowContentType("application/json", "text/plain"))
	h.Router.Post("/api/user/register", h.registerUser)
	h.Router.Post("/api/user/login", h.loginUser)

	h.Router.Route("/api/user", func(r chi.Router) {
		r.Use(jwtauth.Verifier(h.token))
		r.Use(jwtauth.Authenticator(h.token))

		r.Post("/orders", h.loadOrder)
		r.Get("/orders", h.listOrders)
		r.Get("/balance", h.getBalance)
		r.Post("/balance/withdraw", h.withdrawBalance)
		r.Get("/withdrawals", h.listWithdrawals)
	})

	h.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	})

	h.Router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	})
}

func (h *Handler) Shutdown() {
	h.accr.Wait()
	h.storage.Shutdown()
}

func writeText(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write([]byte(message)); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	b, err := json.Marshal(body)
	if err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(b); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}

func writeError(w http.ResponseWriter, err error) {
	code, mes := parsingError(err)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)

	if code != http.StatusNoContent {
		if _, err := w.Write([]byte(mes)); err != nil {
			slog.Error("writeError: ошибка отправки ответа:", slog.Any("error", err))
		}
	}
}

func parsingError(err error) (code int, mes string) {
	var e *custerror.Error

	if errors.As(err, &e) {
		code = e.Code
	} else {
		code = http.StatusInternalServerError
	}

	return code, err.Error()
}

func (h *Handler) getLogin(r *http.Request) (string, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return "", err
	}

	val, ok := claims["login"].(string)
	if !ok {
		return "", fmt.Errorf("токен не содержит информацию о пользователе")
	}

	return val, nil
}
