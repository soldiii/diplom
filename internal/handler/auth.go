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
			user.SupervisorID = ""
			if err != nil {
				NewErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		} else if user.Role == "supervisor" || user.Role == "Supervisor" {
			var supervisor model.Supervisor
			var err error
			id, err = h.services.Authorization.CreateSupervisor(&user, &supervisor)
			user.SupervisorID = ""
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
		/*if err := service.SendEmailAboutRegistration(user.Email); err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		}*/
	}
}

type signInInput struct {
	Email             string `json:"email"`
	EncryptedPassword string `json:"encrypted_password"`
}

func (h *Handler) HandleSignIn() http.HandlerFunc {
	var input signInInput
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		token, err := h.services.Authorization.GenerateToken(input.Email, input.EncryptedPassword)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"token": token})
	}
}

func (h *Handler) HandleGetAllSupervisors() http.HandlerFunc {

	return func(w http.ResponseWriter, _ *http.Request) {

		supervisors, err := h.services.Authorization.GetAllSupervisors()
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(supervisors)
	}
}
