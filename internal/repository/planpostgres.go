package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type PlanPostgres struct {
	db *sqlx.DB
}

func NewPlanPostgres(db *sqlx.DB) *PlanPostgres {
	return &PlanPostgres{db: db}
}

func (r *PlanPostgres) SetPlan(supID int, agentID int) (int, error) {
	var id int
	time := time.Now()
	query := fmt.Sprintf("INSERT INTO %s (supervisor_id, agent_id, internet, tv, convergent, cctv, date_time) VALUES ($1,$2, 0, 0, 0, 0, $3) RETURNING id", plansTable)
	row := r.db.QueryRow(query, supID, agentID, time)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
