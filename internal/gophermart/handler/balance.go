package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
	"github.com/iamamatkazin/diploma-tpl/internal/pkg/custerror"
)

func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request) {
	login, err := h.getLogin(r)
	if err != nil {
		writeError(w, custerror.New(http.StatusUnauthorized, err.Error()))
		return
	}

	balance, err := h.storage.GetBalance(r.Context(), login)
	if err != nil {
		writeError(w, err)
		return
	}

	b, _ := json.Marshal(balance)
	fmt.Println("########2222####", login, string(b))
	writeJSON(w, http.StatusOK, balance)

	/*
		200 — успешная обработка запроса.
		401 — пользователь не авторизован.
		500 — внутренняя ошибка сервера.
	*/
}

func (h *Handler) withdrawBalance(w http.ResponseWriter, r *http.Request) {
	login, err := h.getLogin(r)
	if err != nil {
		writeError(w, custerror.New(http.StatusUnauthorized, err.Error()))
		return
	}

	var withdraw model.Withdraw
	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		writeError(w, custerror.New(http.StatusBadRequest, err.Error()))
		return
	}

	if err := goluhn.Validate(withdraw.Order); err != nil {
		writeError(w, custerror.New(http.StatusUnprocessableEntity, err.Error()))
		return
	}

	code, err := h.storage.WithdrawBalance(r.Context(), login, withdraw)
	if err != nil {
		writeError(w, err)
		return
	}

	switch code {
	case http.StatusOK:
		writeText(w, http.StatusOK, "успешная обработка запроса")
	case http.StatusPaymentRequired:
		writeText(w, http.StatusPaymentRequired, "на счету недостаточно средств")
	default:
		writeText(w, http.StatusInternalServerError, fmt.Sprintf("неизвестный код ошибки: %d", code))
	}

	/*
		200 — успешная обработка запроса;
		401 — пользователь не авторизован;
		402 — на счету недостаточно средств;
		422 — неверный номер заказа;
		500 — внутренняя ошибка сервера.
	*/
}

func (h *Handler) listWithdrawals(w http.ResponseWriter, r *http.Request) {
	login, err := h.getLogin(r)
	if err != nil {
		writeError(w, custerror.New(http.StatusUnauthorized, err.Error()))
		return
	}

	list, err := h.storage.ListWithdrawals(r.Context(), login)
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
		200 — успешная обработка запроса;
		204 — нет ни одного списания.
		401 — пользователь не авторизован.
		500 — внутренняя ошибка сервера.
	*/
}
