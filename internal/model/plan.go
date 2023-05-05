package model

import "time"

type Plan struct {
	ID           int
	SupervisorID string `json:"supervisor_id"`
	AgentID      string `json:"agent_id"`
	Internet     string `json:"internet"`
	TV           string `json:"tv"`
	Convergent   string `json:"Convergent"`
	CCTV         string `json:"cctv"`
	DateTime     time.Time
}
