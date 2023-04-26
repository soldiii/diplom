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
	router.HandleFunc("/info/supervisors", h.HandleGetAllSupervisors()).Methods("GET")
	router.HandleFunc("/auth/sign-in", h.HandleSignIn()).Methods("POST")
	//router.HandleFunc("/users/{id:[0-9]+}", h.HandleGetUserById()).Methods("GET")
	//router.HandleFunc("/wtf", h.HandleWTF()).Methods("GET")
}

//test
/*
func (h *Handler) HandleWTF() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "WATAFACK")
	}
}
*/
