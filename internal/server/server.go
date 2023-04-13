package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/soldiii/diplom/internal/handler"
	"github.com/soldiii/diplom/internal/repository"
)

type Server struct {
	Config     *ServerConfig
	Logger     *logrus.Logger
	Router     *mux.Router
	PostgresDB *repository.PostgresDB
}

func NewServer(cfg *ServerConfig) *Server {

	return &Server{
		Config: cfg,
		Logger: logrus.New(),
		Router: mux.NewRouter(),
	}
}

func (srv *Server) ConfigureLogger() error {
	srv.Logger.SetFormatter(new(logrus.JSONFormatter))
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

	if err := srv.ConfigurePostgresDB(); err != nil {
		return err
	}
	srv.Logger.Info("Сервер запущен")

	return http.ListenAndServe(srv.Config.Address, srv.Router)
}

func (srv *Server) ConfigureRouter() {
	handler := handler.NewHandler()
	handler.InitRoutes(srv.Router)
}

func (srv *Server) ConfigurePostgresDB() error {
	db := repository.NewPostgresDB(srv.Config.DatabaseURL)
	if err := db.OpenPostgresDB(); err != nil {
		return err
	}
	srv.PostgresDB = db
	return nil
}
