package handler

import "net/http"

func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	/*
		200 — пользователь успешно зарегистрирован и аутентифицирован;
		400 — неверный формат запроса;
		409 — логин уже занят;
		500 — внутренняя ошибка сервера.
	*/
}

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	/*
		200 — пользователь успешно зарегистрирован и аутентифицирован;
		400 — неверный формат запроса;
		401 — неверная пара логин/пароль;
		500 — внутренняя ошибка сервера.
	*/
}
