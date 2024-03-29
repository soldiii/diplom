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

func (r *InfoPostgres) GetIsValidByID(uID int) (bool, error) {
	var isValid bool
	query := fmt.Sprintf("SELECT is_valid FROM %s WHERE id = $1", usersTable)
	row := r.db.QueryRow(query, uID)
	if err := row.Scan(&isValid); err != nil {
		return false, err
	}
	return isValid, nil
}

func (r *InfoPostgres) GetSupervisorIDByAgentID(uID int) (int, error) {
	var sup_id int
	query := fmt.Sprintf("SELECT supervisor_id FROM %s WHERE id = $1", agentsTable)
	row := r.db.QueryRow(query, uID)
	if err := row.Scan(&sup_id); err != nil {
		return 0, err
	}
	return sup_id, nil
}

func (r *InfoPostgres) GetFullNameByAgentID(uID int) (string, error) {
	var fullName string
	query := fmt.Sprintf("SELECT CONCAT(surname, ' ', name, COALESCE(CONCAT(' ', patronymic), '')) AS full_name FROM %s u JOIN %s a ON u.id = a.id WHERE a.id = $1", usersTable, agentsTable)
	row := r.db.QueryRow(query, uID)
	if err := row.Scan(&fullName); err != nil {
		return "", err
	}
	return fullName, nil
}

func (r *InfoPostgres) GetSupervisorFullNameByAgentID(uID int) (string, error) {
	var fullName string
	query := fmt.Sprintf("SELECT CONCAT(u.surname, ' ', u.name, ' ', COALESCE(u.patronymic, '')) AS supervisor_name FROM %s u JOIN %s s ON u.id = s.id JOIN %s a ON a.supervisor_id = s.id WHERE a.id = $1", usersTable, supervisorsTable, agentsTable)
	row := r.db.QueryRow(query, uID)
	if err := row.Scan(&fullName); err != nil {
		return "", err
	}
	return fullName, nil
}

func (r *InfoPostgres) GetFullNameBySupID(supID int) (string, error) {
	var fullName string
	query := fmt.Sprintf("SELECT CONCAT(surname, ' ', name, COALESCE(CONCAT(' ', patronymic), '')) AS full_name FROM %s u JOIN %s s ON u.id = s.id WHERE s.id = $1", usersTable, supervisorsTable)
	row := r.db.QueryRow(query, supID)
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

func (r *InfoPostgres) GetReportByAgentID(uID int) (*Rates, error) {
	var flag bool
	report := &Rates{}
	check_query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE agent_id = $1 AND DATE_TRUNC('month', date_time) = DATE_TRUNC('month', CURRENT_DATE)) AS result", reportsTable)
	check_row := r.db.QueryRow(check_query, uID)
	if err := check_row.Scan(&flag); err != nil {
		return nil, err
	}
	if !flag {
		report = &Rates{Internet: "0", TV: "0", Convergent: "0", CCTV: "0"}
		return report, nil
	}
	query := fmt.Sprintf("SELECT SUM(internet), SUM(tv), SUM(convergent), SUM(cctv) FROM %s WHERE agent_id = $1 AND DATE_TRUNC('month', date_time) = DATE_TRUNC('month', CURRENT_DATE)", reportsTable)
	row := r.db.QueryRow(query, uID)
	if err := row.Scan(&report.Internet, &report.TV, &report.Convergent, &report.CCTV); err != nil {
		return nil, err
	}
	return report, nil
}

func (r *InfoPostgres) GetPlanByAgentID(uID int) (*Rates, error) {
	var flag bool
	plan := &Rates{}
	check_query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE agent_id = $1 AND DATE_TRUNC('month', date_time) = DATE_TRUNC('month', CURRENT_DATE)) AS result", plansTable)
	check_row := r.db.QueryRow(check_query, uID)
	if err := check_row.Scan(&flag); err != nil {
		return nil, err
	}
	if !flag {
		plan = &Rates{Internet: "0", TV: "0", Convergent: "0", CCTV: "0"}
		return plan, nil
	}
	query := fmt.Sprintf("SELECT internet, tv, convergent, cctv FROM %s WHERE agent_id = $1 AND DATE_TRUNC('month', date_time) = DATE_TRUNC('month', CURRENT_DATE)", plansTable)
	row := r.db.QueryRow(query, uID)
	if err := row.Scan(&plan.Internet, &plan.TV, &plan.Convergent, &plan.CCTV); err != nil {
		return nil, err
	}
	return plan, nil
}

func (r *InfoPostgres) GetPlanBySupID(supID int) (*Rates, error) {
	var flag bool
	plan := &Rates{}
	check_query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE supervisor_id = $1 AND DATE_TRUNC('month', date_time) = DATE_TRUNC('month', CURRENT_DATE)) AS result", plansTable)
	check_row := r.db.QueryRow(check_query, supID)
	if err := check_row.Scan(&flag); err != nil {
		return nil, err
	}
	if !flag {
		plan = &Rates{Internet: "0", TV: "0", Convergent: "0", CCTV: "0"}
		return plan, nil
	}
	query := fmt.Sprintf("SELECT SUM(p.internet), SUM(p.tv), SUM(p.convergent), SUM(p.cctv) FROM %s p INNER JOIN %s u ON p.agent_id = u.id AND u.is_valid = true WHERE p.supervisor_id = $1 AND DATE_TRUNC('month', p.date_time) = DATE_TRUNC('month', CURRENT_DATE)", plansTable, usersTable)
	row := r.db.QueryRow(query, supID)
	if err := row.Scan(&plan.Internet, &plan.TV, &plan.Convergent, &plan.CCTV); err != nil {
		return nil, err
	}
	return plan, nil
}

func (r *InfoPostgres) CheckForSupervisor(sup_id int) error {
	query := fmt.Sprintf("SELECT id FROM %s WHERE id = $1", supervisorsTable)
	row := r.db.QueryRow(query, sup_id)
	if err := row.Scan(&sup_id); err != nil {
		return err
	}
	return nil
}

type AgentIDAndFullName struct {
	ID       int
	FullName string
}

func (r *InfoPostgres) GetAllAgentsBySupID(supID int) ([]*AgentIDAndFullName, error) {
	query_test := fmt.Sprintf("SELECT COUNT(*) FROM %s a INNER JOIN %s u ON a.id = u.id AND u.is_valid = true WHERE supervisor_id = $1", agentsTable, usersTable)
	rows_test, err := r.db.Query(query_test, supID)
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
		err := errors.New("у супервайзера нет агентов")
		return nil, err
	}

	var agents []*AgentIDAndFullName
	query := fmt.Sprintf("SELECT u.id, CONCAT(surname, ' ', name, COALESCE(CONCAT(' ', patronymic), '')) AS full_name FROM %s u JOIN %s a ON u.id = a.id AND u.is_valid = true WHERE a.supervisor_id = $1", usersTable, agentsTable)
	rows, err := r.db.Query(query, supID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		agent := AgentIDAndFullName{}
		err := rows.Scan(&agent.ID, &agent.FullName)
		if err != nil {
			return nil, err
		}
		agents = append(agents, &agent)
	}
	return agents, nil
}
