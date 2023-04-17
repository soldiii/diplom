package handler

import (
	"encoding/json"
	"net/http"

	"github.com/soldiii/diplom/internal/model"
)

func (h *Handler) HandleSignUp() http.HandlerFunc {
	var user model.User
	var id int

	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if user.Role == "agent" || user.Role == "Agent" {
			var agent model.Agent
			var err error
			id, err = h.services.Authorization.CreateAgent(&user, &agent)
			if err != nil {
				NewErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		} else if user.Role == "supervisor" || user.Role == "Supervisor" {
			var supervisor model.Supervisor
			var err error
			id, err = h.services.Authorization.CreateSupervisor(&user, &supervisor)
			if err != nil {
				NewErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			NewErrorResponse(w, http.StatusInternalServerError, "роль пользователя должна быть либо \"agent\", либо \"supervisor\"")
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
