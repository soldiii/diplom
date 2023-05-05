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
