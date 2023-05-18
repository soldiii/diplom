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
	router.HandleFunc("/refresh", h.HandleRefreshToken()).Methods("POST")
	router.HandleFunc("/info/supervisors", h.HandleGetAllSupervisors()).Methods("GET")
	router.HandleFunc("/info/users/{id:[1-9]+\\d*}/isvalid", h.HandleGetIsValidByID()).Methods("GET")
	router.HandleFunc("/info/agent", h.AuthMiddleware(h.HandleGetInfoAboutAgent())).Methods("GET")
	router.HandleFunc("/info/supervisor", h.AuthMiddleware(h.HandleGetInfoAboutSupervisor())).Methods("GET")
	router.HandleFunc("/info/supervisor/agents", h.AuthMiddleware(h.HandleGetAllAgentsBySupID())).Methods("GET")
	router.HandleFunc("/ad", h.AuthMiddleware(h.HandleGETAndPOSTAd())).Methods("GET", "POST")
	router.HandleFunc("/ad/{id:[1-9]+\\d*}", h.AuthMiddleware(h.HandlePUTAndDELETEAd())).Methods("PUT", "DELETE")
	router.HandleFunc("/report/agent", h.AuthMiddleware(h.HandleGetReportByAgentID())).Methods("GET")
	router.HandleFunc("/report/supervisor", h.AuthMiddleware(h.HandleGetReportBySupervisorID())).Methods("GET")
	router.HandleFunc("/report/supervisor/agents", h.AuthMiddleware(h.HandleGetReportsByAgents())).Methods("GET")
	router.HandleFunc("/report", h.AuthMiddleware(h.HandleCreateReport())).Methods("POST")
	router.HandleFunc("/plan", h.AuthMiddleware(h.HandleGetPlanBySupervisorID())).Methods("GET")
	router.HandleFunc("/plan", h.AuthMiddleware(h.HandleCreatePlan())).Methods("POST")
	router.HandleFunc("/agent/{id:[1-9]+\\d*}", h.AuthMiddleware(h.HandleDeleteAgent())).Methods("DELETE")
}
