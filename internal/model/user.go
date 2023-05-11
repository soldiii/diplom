package model

import "time"

type User struct {
	ID                   int `json:"id" db:"id"`
	Email                string
	Name                 string
	Surname              string
	Patronymic           string
	RegistrationDateTime time.Time
	EncryptedPassword    string
	Role                 string `json:"role" db:"role"`
	SupervisorID         string
	IsValid              bool
}
