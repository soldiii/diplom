package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/soldiii/diplom/internal/handler"
	"github.com/soldiii/diplom/internal/repository"
	"github.com/soldiii/diplom/internal/service"
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

	sqldb, err := srv.ConfigurePostgresDB()
	if err != nil {
		return err
	}
	repo := repository.NewRepository(sqldb)
	srvc := service.NewService(repo)
	srv.ConfigureRouter(srvc)
	srv.Logger.Info("Сервер запущен")

	return http.ListenAndServe(srv.Config.Address, srv.Router)
}

func (srv *Server) ConfigureRouter(service *service.Service) {
	handler := handler.NewHandler(service)
	handler.InitRoutes(srv.Router)
}

func (srv *Server) ConfigurePostgresDB() (*sqlx.DB, error) {
	databaseURL := repository.NewDatabaseURL()
	db := repository.NewPostgresDB(databaseURL)
	sqdb, err := db.OpenPostgresDB()
	if err != nil {
		return nil, err
	}
	srv.PostgresDB = db
	return sqdb, nil
}
