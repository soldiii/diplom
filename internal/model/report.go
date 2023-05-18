package model

import "time"

type Report struct {
	ID         int
	AgentID    int
	Internet   string `json:"internet"`
	TV         string `json:"tv"`
	Convergent string `json:"Convergent"`
	CCTV       string `json:"cctv"`
	DateTime   time.Time
}
