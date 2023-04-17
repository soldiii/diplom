package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/soldiii/diplom/internal/model"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user *model.User) (int, error) {
	var id int
	query_usr := fmt.Sprintf("INSERT INTO %s (email, name, surname, patronymic, reg_date_time, encrypted_password, role, supervisor_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", usersTable)
	row := r.db.QueryRow(query_usr, user.Email, user.Name, user.Surname, user.Patronymic, user.RegistrationDateTime, user.EncryptedPassword, user.Role, user.SupervisorID)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) CreateAgent(user *model.User, agent *model.Agent) (int, error) {
	if _, err := r.CheckForSupervisor(agent.SupervisorID); err != nil {
		err = errors.New("супервайзера с таким id не существует")
		return 0, err
	}
	id, err := r.CreateUser(user)
	if err != nil {
		return 0, err
	}

	agent.ID = id

	query_ag := fmt.Sprintf("INSERT INTO %s (id, supervisor_id) VALUES ($1, $2)", agentsTable)
	r.db.QueryRow(query_ag, agent.ID, agent.SupervisorID)

	return id, nil
}

func (r *AuthPostgres) CreateSupervisor(user *model.User, supervisor *model.Supervisor) (int, error) {
	id, err := r.CreateUser(user)
	if err != nil {
		return 0, err
	}
	supervisor.ID = id

	query_sup := fmt.Sprintf("INSERT INTO %s (id, initials) VALUES ($1, $2)", supervisorsTable)
	r.db.QueryRow(query_sup, supervisor.ID, supervisor.SupervisorInitials)

	return id, nil
}

func (r *AuthPostgres) CheckForSupervisor(sup_id int) (bool, error) {

	query := fmt.Sprintf("SELECT id FROM %s WHERE id = $1", supervisorsTable)
	row := r.db.QueryRow(query, sup_id)
	if err := row.Scan(&sup_id); err != nil {
		return false, err
	}
	return true, nil

}
