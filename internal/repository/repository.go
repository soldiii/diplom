package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/soldiii/diplom/internal/model"
)

type Authorization interface {
	CheckForEmail(email string) error
	CheckForSupervisor(id int) error
	CreateUserTempTable(user *model.UserCode) (int, error)
	GetEmailOfMainSupervisor() (string, error)
	GetSupervisorEmailFromID(id int) (string, error)
	IsDBHaveMainSupervisor() (bool, error)
	CreateMainSupervisor(user *model.User, supervisor *model.Supervisor) (int, error)
	IsRegistrationCodeValid(email string, code string) (bool, error)
	GetAttemptNumber(email string) (int, error)
	MigrateFromTemporaryTable(email string) (int, error)
	GetAllSupervisors() ([]*model.Supervisor, error)
	GetRegistrationTime(email string) (time.Time, error)
	DeleteFromTempTableByEmail(email string)
	IncrementAttemptNumberByEmail(email string)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
