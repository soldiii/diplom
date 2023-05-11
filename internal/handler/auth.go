package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
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
			switch err.Error() {
			case "неверный код", "превышен лимит количества попыток":
				attemptNumber := id
				logrus.Error(err.Error())
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]interface{}{"attempt_number": attemptNumber, "message": err.Error()})
				return
			case "время регистрации истекло":
				NewErrorResponse(w, http.StatusUnauthorized, err.Error())
				return
			default:
				NewErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		/*if err := h.services.Authorization.SetReportAndPlanTables(id); err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}*/
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}

var zeroUser = &model.UserCode{}

func ResetUser(user *model.UserCode) {
	*user = *zeroUser
}

func (h *Handler) HandleSignUp() http.HandlerFunc {
	var user model.UserCode
	return func(w http.ResponseWriter, r *http.Request) {
		ResetUser(&user)
		var id int
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) HandleSignIn() http.HandlerFunc {
	var input signInInput
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		tokens, err := h.services.Authorization.GenerateTokens(input.Email, input.Password)
		if err != nil {
			if err.Error() == "неверный email или пароль" {
				NewErrorResponse(w, http.StatusUnauthorized, err.Error())
				return
			}
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tokens)
	}
}
