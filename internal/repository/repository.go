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
}

type Information interface {
	GetAllSupervisors() ([]*model.Supervisor, error)
	GetUserRoleByID(string) (string, error)
	GetSupervisorIDByAgentID(string) (string, error)
	GetFullNameByAgentID(string) (string, error)
	GetSupervisorFullNameByAgentID(string) (string, error)
	GetReportByAgentID(string) (*Rates, error)
	GetPlanByAgentID(string) (*Rates, error)
	GetFullNameBySupID(string) (string, error)
	GetPlanBySupID(string) (*Rates, error)
	GetAllAgentsBySupID(string) ([]*AgentIDAndFullName, error)
	CheckForSupervisor(string) error
}

type Advertisement interface {
	CreateAd(*model.Advertisement) (int, error)
	UpdateAd(string, string, string) (int, error)
	DeleteAd(string) (int, error)
	GetAdsBySupervisorID(string) ([]*model.Advertisement, error)
	IsSupervisorHaveAds(string) (bool, error)
}

type Report interface {
	SetReport(int) (int, error)
	CreateReport(*model.Report) (int, error)
	IsReportWasCreatedByThisDay(*model.Report) (bool, error)
	UpdateReport(*model.Report) (int, error)
	GetRatesByAgentID(string) (*Rates, error)
	GetRatesBySupervisorIDAndPeriod(string, string) (*Rates, error)
	GetRatesBySupervisorFirstAndLastDates(string, string, string) (*Rates, error)
	GetReportsByAgents(string, string, string) ([]*ReportStructure, error)
}

type Plan interface {
	SetPlan(int, int) (int, error)
	GetPlanBySupervisorID(string) ([]*PlanStructure, error)
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
