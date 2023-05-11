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
	router.HandleFunc("/info/supervisors", h.AuthMiddleware(h.HandleGetAllSupervisors())).Methods("GET")
	router.HandleFunc("/info/users/{id:[1-9]+\\d*}/role", h.HandleGetRoleByID()).Methods("GET")
	router.HandleFunc("/info/users/{id:[1-9]+\\d*}/isvalid", h.HandleGetIsValidByID()).Methods("GET")
	router.HandleFunc("/info/agents/{id:[1-9]+\\d*}", h.HandleGetInfoAboutAgent()).Methods("GET")
	router.HandleFunc("/info/supervisors/{id:[1-9]+\\d*}", h.HandleGetInfoAboutSupervisor()).Methods("GET")
	router.HandleFunc("/info/agents", h.HandleGetAllAgentsBySupID()).Methods("GET")
	router.HandleFunc("/ad", h.HandleGETAndPOSTAd()).Methods("GET", "POST")
	router.HandleFunc("/ad/{id:[1-9]+\\d*}", h.HandlePUTAndDELETEAd()).Methods("PUT", "DELETE")
	router.HandleFunc("/report", h.HandleGetReports()).Methods("GET")
	router.HandleFunc("/report/agents", h.HandleGetReportsByAgents()).Methods("GET")
	router.HandleFunc("/report", h.HandleCreateReport()).Methods("POST")
	router.HandleFunc("/plan", h.HandleGetPlanBySupervisorID()).Methods("GET")
	router.HandleFunc("/plan", h.HandleCreatePlan()).Methods("POST")
	router.HandleFunc("/agent/{id:[1-9]+\\d*}", h.HandleDeleteAgent()).Methods("DELETE")
}
