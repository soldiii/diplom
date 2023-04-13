package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(w http.ResponseWriter, code int, message string) {
	logrus.Error(message)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorResponse{message})
}
