package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/soldiii/diplom/internal/model"
)

type checkCode struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (h *Handler) HandleRegistrationCode() http.HandlerFunc {

	var checkCodeStruct checkCode

	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&checkCodeStruct); err != nil {
			NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		id, err := h.services.Authorization.CompareRegistrationCodes(checkCodeStruct.Email, checkCodeStruct.Code)
		if err != nil {
			if err.Error() == "коды не совпадают" {
				NewErrorResponse(w, http.StatusConflict, err.Error())
				return
			}
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}

func (h *Handler) HandleSignUp() http.HandlerFunc {
	var user model.UserCode
	var id int

	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		switch user.Role {
		case "agent", "Agent":
			var err error
			if user.SupervisorID == "" {
				err := errors.New("необходимо ввести id супервайзера")
				NewErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			id, err = h.services.Authorization.CreateAgent(&user)
			user.SupervisorID = ""
			if err != nil {
				if err.Error() == "почта уже используется" {
					NewErrorResponse(w, http.StatusConflict, err.Error())
					return
				}
				NewErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		case "supervisor", "Supervisor":
			var err error
			if user.SupervisorID != "" {
				err := errors.New("вводить id супервайзера при регистрации не нужно")
				NewErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			id, err = h.services.Authorization.CreateSupervisor(&user)
			user.SupervisorID = ""
			if err != nil {
				if err.Error() == "почта уже используется" {
					NewErrorResponse(w, http.StatusConflict, err.Error())
					return
				}
				NewErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		default:
			NewErrorResponse(w, http.StatusBadRequest, "роль пользователя должна быть либо \"agent\", либо \"supervisor\"")
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
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
			if err.Error() == "в базе данных еще нет супервайзеров" {
				NewErrorResponse(w, http.StatusConflict, err.Error())
				return
			}

			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(supervisors)
	}
}
