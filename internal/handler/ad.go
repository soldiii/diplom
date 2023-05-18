package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soldiii/diplom/internal/model"
)

func (h *Handler) HandleGETAndPOSTAd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			val, err := ParseFromContext(r, user)
			if err != nil {
				NewErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			uID := val.ID
			uRole := val.Role
			ads, err := h.services.Advertisement.GetAdsByUserID(uID, uRole)
			if err != nil {
				if err.Error() == "объявлений нет" {
					NewErrorResponse(w, http.StatusOK, err.Error())
					return
				}
				NewErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ads)
		case "POST":
			var ad model.Advertisement
			val, err := ParseFromContext(r, supervisor)
			if err != nil {
				NewErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			if err := json.NewDecoder(r.Body).Decode(&ad); err != nil {
				NewErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			ad.SupervisorID = val.ID
			id, err := h.services.Advertisement.CreateAd(&ad)
			if err != nil {
				NewErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
		}
	}
}

type PutStructure struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (h *Handler) HandlePUTAndDELETEAd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		adID := vars["id"]
		var id int
		var err error
		switch r.Method {
		case "PUT":
			var ad PutStructure
			if err := json.NewDecoder(r.Body).Decode(&ad); err != nil {
				NewErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			id, err = h.services.Advertisement.UpdateAd(ad.Title, ad.Text, adID)
			if err != nil {
				NewErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		case "DELETE":
			id, err = h.services.Advertisement.DeleteAd(adID)
			if err != nil {
				NewErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	}
}
