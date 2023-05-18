package handler

import (
	"encoding/json"
	"net/http"

	"github.com/soldiii/diplom/internal/model"
)

func (h *Handler) HandleGetPlanBySupervisorID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val, err := ParseFromContext(r, supervisor)
		supID := val.ID
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		plan, err := h.services.Plan.GetPlanBySupervisorID(supID)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(plan)
	}
}

func (h *Handler) HandleCreatePlan() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var plan model.Plan
		val, err := ParseFromContext(r, supervisor)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
			NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		plan.SupervisorID = val.ID
		id, err := h.services.Plan.CreatePlan(&plan)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}
