package service

import (
	"github.com/soldiii/diplom/internal/model"
	"github.com/soldiii/diplom/internal/repository"
)

type Authorization interface {
	CreateAgent(*model.UserCode) (int, error)
	CreateSupervisor(*model.UserCode) (int, error)
	CompareRegistrationCodes(string, string) (int, error)
	GenerateToken(string, string) (string, error)
	IsTimeExpired(string) (bool, error)
	ClearTempTableFromUsersWithExpiredTime() error
	SetReportAndPlanTables(int) error
}

type Information interface {
	GetAllSupervisors() ([]*model.Supervisor, error)
	GetUserRoleByID(string) (string, error)
	GetInfoAboutAgentByID(string) (*InfoAboutAgent, error)
}

type Advertisement interface {
	CreateAd(*model.Advertisement) (int, error)
	UpdateAd(string, string, string) (int, error)
	DeleteAd(string) (int, error)
	GetAdsByUserID(string) ([]*model.Advertisement, error)
}

type Report interface {
	CreateReport(*model.Report) (int, error)
}

type Plan interface {
}
type Service struct {
	Authorization
	Information
	Advertisement
	Report
	Plan
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, repos.Information, repos.Report, repos.Plan),
		Information:   NewInfoService(repos.Information),
		Advertisement: NewAdService(repos.Advertisement, repos.Information),
		Report:        NewReportService(repos.Report, repos.Information),
		Plan:          NewPlanService(repos.Plan, repos.Information),
	}
}
