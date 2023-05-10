package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AgentPostgres struct {
	db *sqlx.DB
}

func NewAgentPostgres(db *sqlx.DB) *AgentPostgres {
	return &AgentPostgres{db: db}
}

func (r *AgentPostgres) DeleteAgent(agentID string) (int, error) {
	var id int
	/*report_query := fmt.Sprintf("DELETE FROM %s where agent_id = $1", reportsTable)
	r.db.QueryRow(report_query, agentID)
	plan_query := fmt.Sprintf("DELETE FROM %s where agent_id = $1", plansTable)
	r.db.QueryRow(plan_query, agentID)
	agent_query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 RETURNING id", agentsTable)
	r.db.QueryRow(agent_query, agentID)
	user_query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 RETURNING id", usersTable)
	row := r.db.QueryRow(user_query, agentID)*/
	query := fmt.Sprintf("UPDATE %s SET is_valid = false where id = $1 RETURNING id", usersTable)
	row := r.db.QueryRow(query, agentID)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
