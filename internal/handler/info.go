package handler

import (
	"encoding/json"
	"net/http"

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

func (h *Handler) HandleGetRoleByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uID := vars["id"]
		role, err := h.services.Information.GetUserRoleByID(uID)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"role": role})
	}
}

func (h *Handler) HandleGetInfoAboutAgent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		agentID := vars["id"]
		info, err := h.services.Information.GetInfoAboutAgentByID(agentID)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(info)
	}
}
