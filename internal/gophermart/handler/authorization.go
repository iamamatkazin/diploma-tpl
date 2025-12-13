package handler

import (
	"encoding/json"
	"net/http"

	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
	"github.com/iamamatkazin/diploma-tpl/internal/pkg/custerror"
)

func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	var login model.Login
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		writeError(w, custerror.New(http.StatusBadRequest, err.Error()))
		return
	}

	if _, err := h.storage.RegisterUser(r.Context(), login); err != nil {
		writeError(w, err)
		return
	}

	_, token, err := h.token.Encode(map[string]any{"login": login.Login})
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	writeText(w, http.StatusOK, "пользователь успешно зарегистрирован и аутентифицирован")

	/*
		200 — пользователь успешно зарегистрирован и аутентифицирован;
		400 — неверный формат запроса;
		409 — логин уже занят;
		500 — внутренняя ошибка сервера.
	*/
}

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {
	var login model.Login
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		writeError(w, custerror.New(http.StatusBadRequest, err.Error()))
		return
	}

	if _, err := h.storage.LoginUser(r.Context(), login); err != nil {
		writeError(w, err)
		return
	}

	_, token, err := h.token.Encode(map[string]any{"login": login.Login})
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	writeText(w, http.StatusOK, "пользователь успешно зарегистрирован и аутентифицирован")
	/*
		200 — пользователь успешно зарегистрирован и аутентифицирован;
		400 — неверный формат запроса;
		401 — неверная пара логин/пароль;
		500 — внутренняя ошибка сервера.
	*/
}
