package handler

import (
	"encoding/json"
	"net/http"

	"github.com/soldiii/diplom/internal/model"
)

func (h *Handler) HandleCreateReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var report model.Report
		val, err := ParseFromContext(r, agent)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
			NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		report.AgentID = val.ID
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

func (h *Handler) HandleGetReportByAgentID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val, err := ParseFromContext(r, agent)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		agentID := val.ID
		data, err := h.services.Report.GetRatesByAgentID(agentID)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}

func (h *Handler) HandleGetReportBySupervisorID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		period := r.URL.Query().Get("period")
		firstDate := r.URL.Query().Get("first_date")
		lastDate := r.URL.Query().Get("last_date")
		val, err := ParseFromContext(r, supervisor)
		supID := val.ID
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
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
	}
}

func (h *Handler) HandleGetReportsByAgents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		firstDate := r.URL.Query().Get("first_date")
		lastDate := r.URL.Query().Get("last_date")
		if firstDate == "" || lastDate == "" {
			NewErrorResponse(w, http.StatusBadRequest, "неверный запрос")
			return
		}
		val, err := ParseFromContext(r, supervisor)
		supID := val.ID
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
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
