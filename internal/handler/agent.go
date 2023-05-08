package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) HandleDeleteAgent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		agentID := vars["id"]
		id, err := h.services.Agent.DeleteAgent(agentID)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}
