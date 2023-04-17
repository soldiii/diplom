package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/soldiii/diplom/internal/model"
)

type Authorization interface {
	CreateAgent(user *model.User, agent *model.Agent) (int, error)
	CreateSupervisor(user *model.User, supervisor *model.Supervisor) (int, error)
	GetAllSupervisors() ([]*model.Supervisor, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
