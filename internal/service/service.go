package service

import (
	"github.com/soldiii/diplom/internal/model"
	"github.com/soldiii/diplom/internal/repository"
)

type Authorization interface {
	CreateAgent(user *model.User, agent *model.Agent) (int, error)
	CreateSupervisor(user *model.User, supervisor *model.Supervisor) (int, error)
	GetAllSupervisors() ([]*model.Supervisor, error)
	GenerateToken(email, password string) (string, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
