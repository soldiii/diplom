package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/soldiii/diplom/internal/model"
	"github.com/soldiii/diplom/internal/repository"
)

type Authorization interface {
	CreateAgent(*model.UserCode) (int, error)
	CreateSupervisor(*model.UserCode) (int, error)
	CompareRegistrationCodes(string, string) (int, error)
	GenerateTokens(string, string) (*Token, error)
	ParseToken(string, bool) (*TokenClaims, error)
	RefreshTokens(string, string, int) (*Token, error)
	CompareRefreshTokens(string, int) (bool, error)
	GenerateTokensByRefresh(string, int) (*Token, error)
	IsTimeExpired(string) (bool, error)
	ClearTempTableFromUsersWithExpiredTime() error
	SetReportAndPlanTables(int) error
	IsTokenExpired(*jwt.NumericDate) bool
}

type Information interface {
	GetAllSupervisors() ([]*model.Supervisor, error)
	GetAllAgentsBySupID(string) ([]*repository.AgentIDAndFullName, error)
	GetUserRoleByID(string) (string, error)
	GetIsValidByID(string) (bool, error)
	GetInfoAboutAgentByID(string) (*InfoAboutAgent, error)
	GetInfoAboutSupervisorByID(string) (*InfoAboutSupervisor, error)
}

type Advertisement interface {
	CreateAd(*model.Advertisement) (int, error)
	UpdateAd(string, string, string) (int, error)
	DeleteAd(string) (int, error)
	GetAdsByUserID(string) ([]*model.Advertisement, error)
}

type Report interface {
	CreateReport(*model.Report) (int, error)
	GetRatesByAgentID(string) (*repository.Rates, error)
	GetRatesBySupervisorIDAndPeriod(string, string) (*repository.Rates, error)
	GetRatesBySupervisorFirstAndLastDates(string, string, string) (*repository.Rates, error)
	GetReportsByAgents(string, string, string) ([]*repository.ReportStructure, error)
}

type Plan interface {
	GetPlanBySupervisorID(string) ([]*repository.PlanStructure, error)
	CreatePlan(*model.Plan) (int, error)
}

type Agent interface {
	DeleteAgent(string) (int, error)
}
type Service struct {
	Authorization
	Information
	Advertisement
	Report
	Plan
	Agent
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, repos.Information, repos.Report, repos.Plan),
		Information:   NewInfoService(repos.Information),
		Advertisement: NewAdService(repos.Advertisement, repos.Information),
		Report:        NewReportService(repos.Report, repos.Information),
		Plan:          NewPlanService(repos.Plan, repos.Information),
		Agent:         NewAgentService(repos.Agent, repos.Information),
	}
}
