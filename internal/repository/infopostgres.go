package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/soldiii/diplom/internal/model"
)

type InfoPostgres struct {
	db *sqlx.DB
}

func NewInfoPostgres(db *sqlx.DB) *InfoPostgres {
	return &InfoPostgres{db: db}
}

func (r *InfoPostgres) GetAllSupervisors() ([]*model.Supervisor, error) {
	query_test := fmt.Sprintf("SELECT COUNT(*) FROM %s", supervisorsTable)
	rows_test, err := r.db.Query(query_test)
	if err != nil {
		return nil, err
	}
	var count int
	for rows_test.Next() {
		if err := rows_test.Scan(&count); err != nil {
			return nil, err
		}
	}
	defer rows_test.Close()
	if count == 0 {
		err := errors.New("в базе данных еще нет супервайзеров")
		return nil, err
	}

	var supervisors []*model.Supervisor
	query := fmt.Sprintf("SELECT * FROM %s", supervisorsTable)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		supervisor := model.Supervisor{}
		err := rows.Scan(&supervisor.ID, &supervisor.SupervisorInitials)
		if err != nil {
			return nil, err
		}
		supervisors = append(supervisors, &supervisor)
	}
	return supervisors, nil
}

func (r *InfoPostgres) GetUserRoleByID(uID string) (string, error) {
	var role string
	query := fmt.Sprintf("SELECT role FROM %s WHERE id = $1", usersTable)
	row := r.db.QueryRow(query, uID)
	if err := row.Scan(&role); err != nil {
		return "", err
	}
	return role, nil
}

func (r *InfoPostgres) GetSupervisorIDByAgentID(uID string) (string, error) {
	var sup_id string
	query := fmt.Sprintf("SELECT supervisor_id FROM %s WHERE id = $1", agentsTable)
	row := r.db.QueryRow(query, uID)
	if err := row.Scan(&sup_id); err != nil {
		return "", err
	}
	return sup_id, nil
}

func (r *InfoPostgres) GetFullNameByID(uID string) (string, error) {
	var fullName string
	query := fmt.Sprintf("SELECT CONCAT(surname, ' ', name, COALESCE(CONCAT(' ', patronymic), '')) AS full_name FROM %s u JOIN %s a ON u.id = a.id WHERE a.id = $1", usersTable, agentsTable)
	row := r.db.QueryRow(query, uID)
	if err := row.Scan(&fullName); err != nil {
		return "", err
	}
	return fullName, nil
}

func (r *InfoPostgres) GetSupervisorFullNameByID(uID string) (string, error) {
	var fullName string
	query := fmt.Sprintf("SELECT CONCAT(u.surname, ' ', u.name, ' ', COALESCE(u.patronymic, '')) AS supervisor_name FROM %s u JOIN %s s ON u.id = s.id JOIN %s a ON a.supervisor_id = s.id WHERE a.id = $1", usersTable, supervisorsTable, agentsTable)
	row := r.db.QueryRow(query, uID)
	if err := row.Scan(&fullName); err != nil {
		return "", err
	}
	return fullName, nil
}

type Rates struct {
	Internet   string
	TV         string
	Convergent string
	CCTV       string
}

func (r *InfoPostgres) GetReportByID(uID string) (*Rates, error) {
	query := fmt.Sprintf("SELECT SUM(internet), SUM(tv), SUM(convergent), SUM(cctv) FROM %s WHERE agent_id = $1 AND DATE_TRUNC('month', date_time) = DATE_TRUNC('month', CURRENT_DATE)", reportsTable)
	row := r.db.QueryRow(query, uID)
	report := &Rates{}
	if err := row.Scan(&report.Internet, &report.TV, &report.Convergent, &report.CCTV); err != nil {
		return nil, err
	}
	return report, nil
}

func (r *InfoPostgres) GetPlanByID(uID string) (*Rates, error) {
	query := fmt.Sprintf("SELECT internet, tv, convergent, cctv FROM %s WHERE agent_id = $1 AND DATE_TRUNC('month', date_time) = DATE_TRUNC('month', CURRENT_DATE)", plansTable)
	row := r.db.QueryRow(query, uID)
	plan := &Rates{}
	if err := row.Scan(&plan.Internet, &plan.TV, &plan.Convergent, &plan.CCTV); err != nil {
		return nil, err
	}
	return plan, nil
}
