package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/soldiii/diplom/internal/model"
)

type PlanPostgres struct {
	db *sqlx.DB
}

func NewPlanPostgres(db *sqlx.DB) *PlanPostgres {
	return &PlanPostgres{db: db}
}

type PlanStructure struct {
	ID         int
	FullName   string
	Internet   int
	TV         int
	Convergent int
	CCTV       int
}

func (r *PlanPostgres) GetPlanBySupervisorID(supID string) ([]*PlanStructure, error) {
	var plans []*PlanStructure
	query := fmt.Sprintf("SELECT a.id, CONCAT(u.surname, ' ', u.name, COALESCE(CONCAT(' ', u.patronymic), '')) AS full_name, COALESCE (SUM(CASE WHEN DATE_TRUNC('month', p.date_time) = DATE_TRUNC('month', CURRENT_DATE) THEN p.internet ELSE 0 END), 0) AS internet, COALESCE (SUM(CASE WHEN DATE_TRUNC('month', p.date_time) = DATE_TRUNC('month', CURRENT_DATE) THEN p.tv ELSE 0 END), 0) AS tv, COALESCE (SUM(CASE WHEN DATE_TRUNC('month', p.date_time) = DATE_TRUNC('month', CURRENT_DATE) THEN p.convergent ELSE 0 END), 0) AS convergent, COALESCE (SUM(CASE WHEN DATE_TRUNC('month', p.date_time) = DATE_TRUNC('month', CURRENT_DATE) THEN p.cctv ELSE 0 END), 0) AS cctv FROM %s a LEFT JOIN %s p ON a.id = p.agent_id JOIN %s u ON a.id = u.id AND u.is_valid = true WHERE a.supervisor_id = $1 GROUP BY a.id, u.surname, u.name, u.patronymic ORDER BY a.id", agentsTable, plansTable, usersTable)
	rows, err := r.db.Query(query, supID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		plan := PlanStructure{}
		err := rows.Scan(&plan.ID, &plan.FullName, &plan.Internet, &plan.TV, &plan.Convergent, &plan.CCTV)
		if err != nil {
			return nil, err
		}
		plans = append(plans, &plan)
	}
	return plans, nil
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

func (r *PlanPostgres) IsPlanWasCreatedByThisMonth(plan *model.Plan) (bool, error) {
	var flag bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE agent_id = $1 AND DATE_TRUNC('month', date_time) = DATE_TRUNC('month', CURRENT_DATE)) AS result", plansTable)
	row := r.db.QueryRow(query, plan.AgentID)
	if err := row.Scan(&flag); err != nil {
		return false, err
	}
	if flag {
		return true, nil
	}
	return false, nil
}

func (r *PlanPostgres) UpdatePlan(plan *model.Plan) (int, error) {
	var id int
	query := fmt.Sprintf("UPDATE %s SET internet = $3, tv = $4, convergent = $5, cctv = $6 WHERE supervisor_id = $1 AND agent_id = $2 AND DATE_TRUNC('month', date_time) = DATE_TRUNC('month', CURRENT_DATE) RETURNING id", plansTable)
	row := r.db.QueryRow(query, plan.SupervisorID, plan.AgentID, plan.Internet, plan.TV, plan.Convergent, plan.CCTV)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PlanPostgres) CreatePlan(plan *model.Plan) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (supervisor_id, agent_id, internet, tv, convergent, cctv, date_time) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", plansTable)
	row := r.db.QueryRow(query, plan.SupervisorID, plan.AgentID, plan.Internet, plan.TV, plan.Convergent, plan.CCTV, plan.DateTime)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
