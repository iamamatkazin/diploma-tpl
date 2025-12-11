package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	"github.com/iamamatkazin/diploma-tpl/internal/config"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/repository"
	"github.com/iamamatkazin/diploma-tpl/internal/pkg/custerror"
)

type Handler struct {
	token   *jwtauth.JWTAuth
	storage repository.Storager
	Router  *chi.Mux
	cfg     *config.Config
}

func New(ctx context.Context, cfg *config.Config) (*Handler, error) {
	storage, err := repository.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		storage: storage,
		cfg:     cfg,
		token:   jwtauth.New("HS256", []byte(cfg.SecretKey), nil),
	}

	h.Router = chi.NewRouter()
	h.listRoute()

	return h, nil
}

func (h *Handler) listRoute() {
	h.Router.Use(middleware.AllowContentType("application/json", "text/plain"))
	h.Router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.registerUser)
		r.Post("/login", h.loginUser)

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

	if _, err := w.Write([]byte(mes)); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
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
