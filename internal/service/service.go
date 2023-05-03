package service

import (
	"github.com/soldiii/diplom/internal/model"
	"github.com/soldiii/diplom/internal/repository"
)

type Authorization interface {
	CreateAgent(user *model.UserCode) (int, error)
	CreateSupervisor(user *model.UserCode) (int, error)
	GetAllSupervisors() ([]*model.Supervisor, error)
	CompareRegistrationCodes(email string, code string) (int, error)
	GenerateToken(email, password string) (string, error)
	IsTimeExpired(email string) (bool, error)
	ClearTempTableFromUsersWithExpiredTime() error
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{Authorization: NewAuthService(repos.Authorization)}
}
