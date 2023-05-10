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
		id, err := h.services.Report.CreateReport(&report)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}

func (h *Handler) HandleGetReports() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		supID := r.URL.Query().Get("supervisor_id")
		agentID := r.URL.Query().Get("agent_id")
		period := r.URL.Query().Get("period")
		firstDate := r.URL.Query().Get("first_date")
		lastDate := r.URL.Query().Get("last_date")

		if agentID != "" && supID == "" && period == "" && firstDate == "" && lastDate == "" {
			data, err := h.services.Report.GetRatesByAgentID(agentID)
			if err != nil {
				NewErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(data)
			return
		}

		if supID != "" && agentID == "" {
			if period != "" && firstDate == "" && lastDate == "" {
				if period != "month" && period != "day" && period != "week" {
					NewErrorResponse(w, http.StatusBadRequest, "неверное значение периода")
					return
				}
				data, err := h.services.Report.GetRatesBySupervisorIDAndPeriod(supID, period)
				if err != nil {
					NewErrorResponse(w, http.StatusInternalServerError, err.Error())
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(data)
				return
			}
			if firstDate != "" && lastDate != "" && period == "" {
				data, err := h.services.Report.GetRatesBySupervisorFirstAndLastDates(supID, firstDate, lastDate)
				if err != nil {
					NewErrorResponse(w, http.StatusInternalServerError, err.Error())
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(data)
				return
			}
			NewErrorResponse(w, http.StatusBadRequest, "неверное тело запроса")
			return
		}
	}
}

func (h *Handler) HandleGetReportsByAgents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		supID := r.URL.Query().Get("supervisor_id")
		firstDate := r.URL.Query().Get("first_date")
		lastDate := r.URL.Query().Get("last_date")
		if supID == "" || firstDate == "" || lastDate == "" {
			NewErrorResponse(w, http.StatusBadRequest, "неверный запрос")
			return
		}
		data, err := h.services.Report.GetReportsByAgents(supID, firstDate, lastDate)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}
