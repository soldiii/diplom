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
	GetAllAgentsBySupID(int) ([]*repository.AgentIDAndFullName, error)
	GetIsValidByID(int) (bool, error)
	GetInfoAboutAgentByID(int) (*InfoAboutAgent, error)
	GetInfoAboutSupervisorByID(int) (*InfoAboutSupervisor, error)
}

type Advertisement interface {
	CreateAd(*model.Advertisement) (int, error)
	UpdateAd(string, string, string) (int, error)
	DeleteAd(string) (int, error)
	GetAdsByUserID(int, string) ([]*model.Advertisement, error)
}

type Report interface {
	CreateReport(*model.Report) (int, error)
	GetRatesByAgentID(int) (*repository.Rates, error)
	GetRatesBySupervisorIDAndPeriod(int, string) (*repository.Rates, error)
	GetRatesBySupervisorFirstAndLastDates(int, string, string) (*repository.Rates, error)
	GetReportsByAgents(int, string, string) ([]*repository.ReportStructure, error)
}

type Plan interface {
	GetPlanBySupervisorID(int) ([]*repository.PlanStructure, error)
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
