package model

type Advertisement struct {
	ID           int
	SupervisorID int
	Title        string `json:"title"`
	Text         string `json:"text"`
}
