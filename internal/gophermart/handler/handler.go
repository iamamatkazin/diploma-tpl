package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/iamamatkazin/diploma-tpl/internal/config"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/repository"
)

type Handler struct {
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
	}

	h.Router = chi.NewRouter()
	h.listRoute()

	return h, nil
}

func (h *Handler) listRoute() {
	h.Router.With(middleware.AllowContentType("application/json")).Post("/api/user/register", h.registerUser)
	h.Router.With(middleware.AllowContentType("application/json")).Post("/api/user/login", h.loginUser)
	h.Router.With(middleware.AllowContentType("text/plain")).Post("/api/user/orders", h.loadOrder)
	h.Router.With(middleware.AllowContentType("application/json")).Get("/api/user/orders", h.listOrders)
	h.Router.With(middleware.AllowContentType("application/json")).Get("/api/user/balance", h.getBalance)
	h.Router.With(middleware.AllowContentType("application/json")).Post("/api/user/balance/withdraw", h.withdrawBalance)
	h.Router.With(middleware.AllowContentType("application/json")).Get("/api/user/withdrawals", h.listWithdrawals)

	h.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	})

	h.Router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	})
}

func writeText(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write([]byte(message)); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}

func writeJSON(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(body); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}

func writeHTML(w http.ResponseWriter, status int, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write([]byte(html)); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}

func (h *Handler) Shutdown() {
	// h.storage.Shutdown()
}
