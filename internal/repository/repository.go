package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/soldiii/diplom/internal/model"
)

type Authorization interface {
	CheckForEmail(string) error
	IsTempTableHaveUser(string) (bool, error)
	GetAttemptNumberByEmail(string) (int, error)
	GetCodeByEmail(string) (string, error)
	CreateUserTempTable(*model.UserCode) (int, error)
	GetEmailOfMainSupervisor() (string, error)
	GetSupervisorEmailFromID(int) (string, error)
	IsDBHaveMainSupervisor() (bool, error)
	CreateMainSupervisor(*model.User, *model.Supervisor) (int, error)
	IsRegistrationCodeValid(string, string) (bool, error)
	MigrateFromTemporaryTable(string) (int, error)
	GetRegistrationTimeByEmail(string) (time.Time, error)
	GetUsersEmailsWithExpiredTime(time.Time, int64) ([]string, error)
	IncrementAttemptNumberByEmail(string)
	DeleteFromTempTableByEmail(string)
	GetUser(string, string) (*model.User, error)
	GetPassword(string) (string, error)
	IsEmailValid(string) (bool, error)
	PostRefreshToken(int, string) (int, error)
	IsUserHaveRefreshToken(int) (bool, error)
	UpdateRefreshToken(int, string) (int, error)
	CompareRefreshTokens(string, int) (bool, error)
}

type Information interface {
	GetAllSupervisors() ([]*model.Supervisor, error)
	GetIsValidByID(int) (bool, error)
	GetSupervisorIDByAgentID(int) (int, error)
	GetFullNameByAgentID(int) (string, error)
	GetSupervisorFullNameByAgentID(int) (string, error)
	GetReportByAgentID(int) (*Rates, error)
	GetPlanByAgentID(int) (*Rates, error)
	GetFullNameBySupID(int) (string, error)
	GetPlanBySupID(int) (*Rates, error)
	GetAllAgentsBySupID(int) ([]*AgentIDAndFullName, error)
	CheckForSupervisor(int) error
}

type Advertisement interface {
	CreateAd(*model.Advertisement) (int, error)
	UpdateAd(string, string, string) (int, error)
	DeleteAd(string) (int, error)
	GetAdsBySupervisorID(int) ([]*model.Advertisement, error)
	IsSupervisorHaveAds(int) (bool, error)
}

type Report interface {
	SetReport(int) (int, error)
	CreateReport(*model.Report) (int, error)
	IsReportWasCreatedByThisDay(*model.Report) (bool, error)
	UpdateReport(*model.Report) (int, error)
	GetRatesByAgentID(int) (*Rates, error)
	GetRatesBySupervisorIDAndPeriod(int, string) (*Rates, error)
	GetRatesBySupervisorFirstAndLastDates(int, string, string) (*Rates, error)
	GetReportsByAgents(int, string, string) ([]*ReportStructure, error)
}

type Plan interface {
	SetPlan(int, int) (int, error)
	GetPlanBySupervisorID(int) ([]*PlanStructure, error)
	IsPlanWasCreatedByThisMonth(*model.Plan) (bool, error)
	UpdatePlan(*model.Plan) (int, error)
	CreatePlan(*model.Plan) (int, error)
}

type Agent interface {
	DeleteAgent(string) (int, error)
}

type Repository struct {
	Authorization
	Information
	Advertisement
	Report
	Plan
	Agent
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Information:   NewInfoPostgres(db),
		Advertisement: NewAdPostgres(db),
		Report:        NewReportPostgres(db),
		Plan:          NewPlanPostgres(db),
		Agent:         NewAgentPostgres(db),
	}
}
