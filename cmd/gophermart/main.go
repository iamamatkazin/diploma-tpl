package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamamatkazin/diploma-tpl/internal/config"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/handler"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		slog.Error("Ошибка чтения конфигурации:", slog.Any("error", err))
		os.Exit(2)
	}

	handler, err := handler.New(ctx, cfg)
	if err != nil {
		slog.Error("Ошибка создания сервера:", slog.Any("error", err))
		os.Exit(2)
	}

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: handler.Router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	exit := make(chan struct{})

	go func() {
		slog.Info("Запуск сервера")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Ошибка запуска сервера:", slog.Any("error", err))
			close(exit)
		}
	}()

	select {
	case <-quit:
		cancel()
		handler.Shutdown()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Ошибка остановки сервера:", slog.Any("error", err))
		}

	case <-exit:
		os.Exit(2)
	}

	slog.Info("Выключение сервера")
}
