package handler

import (
	"encoding/json"
	"net/http"

	"github.com/soldiii/diplom/internal/model"
)

func (h *Handler) HandleSignUp() http.HandlerFunc {
	var user model.User

	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		id, err := h.services.Authorization.CreateUser(&user)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}

func (h *Handler) HandleSignIn() http.HandlerFunc {

	return nil
}

func (h *Handler) HandleGetUserById() http.HandlerFunc {

	return nil
}
