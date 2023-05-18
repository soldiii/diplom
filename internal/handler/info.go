package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) HandleGetAllSupervisors() http.HandlerFunc {

	return func(w http.ResponseWriter, _ *http.Request) {

		supervisors, err := h.services.Information.GetAllSupervisors()
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

func (h *Handler) HandleGetInfoAboutAgent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val, err := ParseFromContext(r, agent)
		uID := val.ID
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		info, err := h.services.Information.GetInfoAboutAgentByID(uID)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(info)
	}
}

func (h *Handler) HandleGetAllAgentsBySupID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val, err := ParseFromContext(r, supervisor)
		supID := val.ID
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		agents, err := h.services.Information.GetAllAgentsBySupID(supID)
		if err != nil {
			if err.Error() == "у супервайзера нет агентов" {
				NewErrorResponse(w, http.StatusOK, err.Error())
				return
			}
			if err.Error() == "супервайзер с таким id не существует" {
				NewErrorResponse(w, http.StatusConflict, err.Error())
				return
			}

			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(agents)
	}
}

func (h *Handler) HandleGetInfoAboutSupervisor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val, err := ParseFromContext(r, supervisor)
		supID := val.ID
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		info, err := h.services.Information.GetInfoAboutSupervisorByID(supID)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(info)
	}
}

func (h *Handler) HandleGetIsValidByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		agentID := vars["id"]
		ID, err := strconv.Atoi(agentID)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		isValid, err := h.services.Information.GetIsValidByID(ID)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"IsValid": isValid})
	}
}
