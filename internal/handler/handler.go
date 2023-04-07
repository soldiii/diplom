package handler

import "github.com/gorilla/mux"

type Handler struct {
}

func (h *Handler) InitRoutes(router *mux.Router) {

	router.HandleFunc("/auth/sign-up/", h.HandleSignUp()).Methods("POST")
	router.HandleFunc("/auth/sign-up/", h.HandleSignIn()).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}/", h.HandleGetUserById()).Methods("GET")
}
