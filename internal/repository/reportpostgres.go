package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/soldiii/diplom/internal/model"
)

type ReportPostgres struct {
	db *sqlx.DB
}

func NewReportPostgres(db *sqlx.DB) *ReportPostgres {
	return &ReportPostgres{db: db}
}

func (r *ReportPostgres) CreateReport(report *model.Report) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (agent_id, internet, tv, convergent, cctv, date_time) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", reportsTable)
	row := r.db.QueryRow(query, report.AgentID, report.Internet, report.TV, report.Convergent, report.CCTV, report.DateTime)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ReportPostgres) SetReport(agentID int) (int, error) {
	var id int
	time := time.Now()
	query := fmt.Sprintf("INSERT INTO %s (agent_id, internet, tv, convergent, cctv, date_time) VALUES ($1, 0, 0, 0, 0, $2) RETURNING id", reportsTable)
	row := r.db.QueryRow(query, agentID, time)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ReportPostgres) IsReportWasCreatedByThisDay(report *model.Report) (bool, error) {
	var flag bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE agent_id = $1 AND DATE_TRUNC('day', date_time) = DATE_TRUNC('day', CURRENT_DATE)) AS result", reportsTable)
	row := r.db.QueryRow(query, report.AgentID)
	if err := row.Scan(&flag); err != nil {
		return false, err
	}
	if flag {
		return true, nil
	}
	return false, nil
}

func (r *ReportPostgres) UpdateReport(report *model.Report) (int, error) {
	var id int
	query := fmt.Sprintf("UPDATE %s SET internet = internet + $2, tv = tv + $3, convergent = convergent + $4, cctv = cctv + $5 WHERE agent_id = $1 AND DATE_TRUNC('day', date_time) = DATE_TRUNC('day', CURRENT_DATE) RETURNING id", reportsTable)
	row := r.db.QueryRow(query, report.AgentID, report.Internet, report.TV, report.Convergent, report.CCTV)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ReportPostgres) GetRatesByAgentID(agentID string) (*Rates, error) {
	var flag bool
	report := &Rates{}
	check_query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE agent_id = $1 AND DATE_TRUNC('day', date_time) = DATE_TRUNC('day', CURRENT_DATE)) AS result", reportsTable)
	check_row := r.db.QueryRow(check_query, agentID)
	if err := check_row.Scan(&flag); err != nil {
		return nil, err
	}
	if !flag {
		report = &Rates{Internet: "0", TV: "0", Convergent: "0", CCTV: "0"}
		return report, nil
	}
	query := fmt.Sprintf("SELECT internet, tv, convergent, cctv FROM %s WHERE agent_id = $1 AND DATE_TRUNC('day', date_time) = DATE_TRUNC('day', CURRENT_DATE)", reportsTable)
	row := r.db.QueryRow(query, agentID)
	if err := row.Scan(&report.Internet, &report.TV, &report.Convergent, &report.CCTV); err != nil {
		return nil, err
	}
	return report, nil
}

func (r *ReportPostgres) GetRatesBySupervisorIDAndPeriod(supID, period string) (*Rates, error) {
	var flag bool
	report := &Rates{}
	check_query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s r JOIN %s a ON r.agent_id = a.id WHERE a.supervisor_id = $1 AND DATE_TRUNC($2, date_time) = DATE_TRUNC($2, CURRENT_DATE)) AS result", reportsTable, agentsTable)
	check_row := r.db.QueryRow(check_query, supID, period)
	if err := check_row.Scan(&flag); err != nil {
		return nil, err
	}
	if !flag {
		report = &Rates{Internet: "0", TV: "0", Convergent: "0", CCTV: "0"}
		return report, nil
	}

	query := fmt.Sprintf("SELECT SUM(internet), SUM(tv), SUM(convergent), SUM(cctv) FROM %s r JOIN %s a ON r.agent_id = a.id WHERE a.supervisor_id = $1 AND DATE_TRUNC($2, date_time) = DATE_TRUNC($2, CURRENT_DATE)", reportsTable, agentsTable)
	row := r.db.QueryRow(query, supID, period)
	if err := row.Scan(&report.Internet, &report.TV, &report.Convergent, &report.CCTV); err != nil {
		return nil, err
	}
	return report, nil
}

func (r *ReportPostgres) GetRatesBySupervisorFirstAndLastDates(supID, firstDate, lastDate string) (*Rates, error) {
	var flag bool
	report := &Rates{}
	check_query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s r JOIN %s a ON r.agent_id = a.id WHERE a.supervisor_id = $1 AND DATE_TRUNC('day', date_time) >= $2::timestamp AND date_time <= $3::timestamp) AS result", reportsTable, agentsTable)
	check_row := r.db.QueryRow(check_query, supID, firstDate, lastDate)
	if err := check_row.Scan(&flag); err != nil {
		return nil, err
	}
	if !flag {
		report = &Rates{Internet: "0", TV: "0", Convergent: "0", CCTV: "0"}
		return report, nil
	}

	query := fmt.Sprintf("SELECT SUM(internet), SUM(tv), SUM(convergent), SUM(cctv) FROM %s r JOIN %s a ON r.agent_id = a.id WHERE a.supervisor_id = $1 AND DATE_TRUNC('day', date_time) >= $2::timestamp AND date_time <= $3::timestamp", reportsTable, agentsTable)
	row := r.db.QueryRow(query, supID, firstDate, lastDate)
	if err := row.Scan(&report.Internet, &report.TV, &report.Convergent, &report.CCTV); err != nil {
		return nil, err
	}
	return report, nil
}

type ReportStructure struct {
	ID         int
	FullName   string
	Internet   int
	TV         int
	Convergent int
	CCTV       int
}

func (r *ReportPostgres) GetReportsByAgents(supID, firstDate, lastDate string) ([]*ReportStructure, error) {
	var reports []*ReportStructure
	query := fmt.Sprintf("SELECT a.id, CONCAT(u.surname, ' ', u.name, COALESCE(CONCAT(' ', u.patronymic), '')) AS full_name, COALESCE (SUM(CASE WHEN DATE_TRUNC('day', r.date_time) >= $2::timestamp AND r.date_time <= $3::timestamp THEN r.internet ELSE 0 END), 0) AS internet, COALESCE (SUM(CASE WHEN DATE_TRUNC('day', r.date_time) >= $2::timestamp AND r.date_time <= $3::timestamp THEN r.tv ELSE 0 END), 0) AS tv, COALESCE (SUM(CASE WHEN DATE_TRUNC('day', r.date_time) >= $2::timestamp AND r.date_time <= $3::timestamp THEN r.convergent ELSE 0 END), 0) AS convergent, COALESCE (SUM(CASE WHEN DATE_TRUNC('day', r.date_time) >= $2::timestamp AND r.date_time <= $3::timestamp THEN r.cctv ELSE 0 END), 0) AS cctv FROM %s a LEFT JOIN %s r ON a.id = r.agent_id JOIN %s u ON a.id = u.id WHERE a.supervisor_id = $1 GROUP BY a.id, u.surname, u.name, u.patronymic ORDER BY a.id", agentsTable, reportsTable, usersTable)
	rows, err := r.db.Query(query, supID, firstDate, lastDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		report := ReportStructure{}
		err := rows.Scan(&report.ID, &report.FullName, &report.Internet, &report.TV, &report.Convergent, &report.CCTV)
		if err != nil {
			return nil, err
		}
		reports = append(reports, &report)
	}
	return reports, nil
}
