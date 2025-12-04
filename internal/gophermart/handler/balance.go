package handler

import "net/http"

func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request) {
	/*
		200 — успешная обработка запроса.
		401 — пользователь не авторизован.
		500 — внутренняя ошибка сервера.
	*/
}

func (h *Handler) withdrawBalance(w http.ResponseWriter, r *http.Request) {
	/*
		200 — успешная обработка запроса;
		401 — пользователь не авторизован;
		402 — на счету недостаточно средств;
		422 — неверный номер заказа;
		500 — внутренняя ошибка сервера.
	*/
}

func (h *Handler) listWithdrawals(w http.ResponseWriter, r *http.Request) {
	/*
		200 — успешная обработка запроса;
		204 — нет ни одного списания.
		401 — пользователь не авторизован.
		500 — внутренняя ошибка сервера.
	*/
}
