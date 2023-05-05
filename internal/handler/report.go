package handler

import (
	"encoding/json"
	"net/http"

	"github.com/soldiii/diplom/internal/model"
)

func (h *Handler) HandleCreateReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var report model.Report
		if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
			NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		id, err := h.services.CreateReport(&report)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}
