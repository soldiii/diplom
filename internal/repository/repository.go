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
	CheckForSupervisor(int) error
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
	GetFullNameByID(string) (string, error)
	GetSupervisorFullNameByID(string) (string, error)
	GetReportByID(string) (*Rates, error)
	GetPlanByID(string) (*Rates, error)
}

type Advertisement interface {
	CreateAd(*model.Advertisement) (int, error)
	UpdateAd(string, string, string) (int, error)
	DeleteAd(string) (int, error)
	GetAdsBySupervisorID(string) ([]*model.Advertisement, error)
}

type Report interface {
	SetReport(int) (int, error)
	CreateReport(*model.Report) (int, error)
}

type Plan interface {
	SetPlan(int, int) (int, error)
}

type Repository struct {
	Authorization
	Information
	Advertisement
	Report
	Plan
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Information:   NewInfoPostgres(db),
		Advertisement: NewAdPostgres(db),
		Report:        NewReportPostgres(db),
		Plan:          NewPlanPostgres(db),
	}
}
