package handler

import (
	"github.com/gorilla/mux"
	"github.com/soldiii/diplom/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes(router *mux.Router) {

	router.HandleFunc("/auth/sign-up", h.HandleSignUp()).Methods("POST")
	router.HandleFunc("/auth/sign-up/check-up", h.HandleRegistrationCode()).Methods("POST")
	router.HandleFunc("/auth/sign-in", h.HandleSignIn()).Methods("POST")
	router.HandleFunc("/info/supervisors", h.HandleGetAllSupervisors()).Methods("GET")
	router.HandleFunc("/info/users/{id:[1-9]+\\d*}/role", h.HandleGetRoleByID()).Methods("GET")
	router.HandleFunc("/info/agent/{id:[1-9]+\\d*}", h.HandleGetInfoAboutAgent()).Methods("GET")
	router.HandleFunc("/ad", h.HandleGETAndPOSTAd()).Methods("GET", "POST")
	router.HandleFunc("/ad/{id:[1-9]+\\d*}", h.HandlePUTAndDELETEAd()).Methods("PUT", "DELETE")
	router.HandleFunc("/report", h.HandleCreateReport()).Methods("POST")

	//router.HandleFunc("/users/{id:[0-9]+}", h.HandleGetUserById()).Methods("GET")
}
