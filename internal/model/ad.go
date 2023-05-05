package model

type Advertisement struct {
	ID           int
	SupervisorID string `json:"supervisor_id"`
	Title        string `json:"title"`
	Text         string `json:"text"`
}
