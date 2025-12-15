package handler

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"

	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
	"github.com/iamamatkazin/diploma-tpl/internal/pkg/custerror"
)

func (h *Handler) loadOrder(w http.ResponseWriter, r *http.Request) {
	login, err := h.getLogin(r)
	if err != nil {
		writeError(w, custerror.New(http.StatusUnauthorized, err.Error()))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, custerror.New(http.StatusBadRequest, err.Error()))
		return
	}

	order := string(body)
	if err := goluhn.Validate(order); err != nil {
		writeError(w, custerror.New(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	currentLogin, err := h.storage.LoadOrder(r.Context(), login, order)
	if err != nil {
		writeError(w, err)
		return
	}

	switch {
	case currentLogin == "":
		writeText(w, http.StatusAccepted, "новый номер заказа принят в обработку")

		select {
		case h.chOrder <- model.UserOrder{Login: login, Order: order}:
		default:
			slog.Info("Занят канал отправки заказа в систему расчета начислений")
		}

	case currentLogin == login:
		writeText(w, http.StatusOK, "номер заказа уже был загружен этим пользователем")

	case currentLogin != login:
		writeText(w, http.StatusConflict, "номер заказа уже был загружен другим пользователем")

	default:
	}

	/*
		200 — номер заказа уже был загружен этим пользователем;
		202 — новый номер заказа принят в обработку;
		400 — неверный формат запроса;
		401 — пользователь не аутентифицирован;
		409 — номер заказа уже был загружен другим пользователем;
		422 — неверный формат номера заказа;
		500 — внутренняя ошибка сервера.
	*/
}

func (h *Handler) listOrders(w http.ResponseWriter, r *http.Request) {
	login, err := h.getLogin(r)
	if err != nil {
		writeError(w, custerror.New(http.StatusUnauthorized, err.Error()))
		return
	}

	list, err := h.storage.ListOrders(r.Context(), login)
	if err != nil {
		writeError(w, err)
		return
	}

	if len(list) == 0 {
		writeError(w, custerror.New(http.StatusNoContent, ""))
		return
	}

	writeJSON(w, http.StatusOK, list)

	/*
		200 — успешная обработка запроса
		204 — нет данных для ответа.
		401 — пользователь не авторизован.
		500 — внутренняя ошибка сервера.
	*/
}
