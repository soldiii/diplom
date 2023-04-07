package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/soldiii/diplom/internal/handler"
)

type Server struct {
	Config *ServerConfig
	Logger *logrus.Logger
	Router *mux.Router
}

func NewServer(cfg *ServerConfig) *Server {

	return &Server{
		Config: cfg,
		Logger: logrus.New(),
		Router: mux.NewRouter(),
	}
}

func (srv *Server) ConfigureLogger() error {
	srv.Logger.SetFormatter(&logrus.JSONFormatter{})
	lvl, err := logrus.ParseLevel(srv.Config.LogLevel)
	if err != nil {
		return err
	}
	srv.Logger.SetLevel(lvl)
	return nil
}

func (srv *Server) RunServer() error {

	if err := srv.ConfigureLogger(); err != nil {
		return err
	}

	srv.ConfigureRouter()

	srv.Logger.Info("Сервер запущен")

	return http.ListenAndServe(srv.Config.Addr, srv.Router)
}

func (srv *Server) ConfigureRouter() {
	handler := handler.NewHandler()
	handler.InitRoutes(srv.Router)
}
