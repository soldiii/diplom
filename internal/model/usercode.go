package model

import "time"

type UserCode struct {
	ID                   int
	Email                string `json:"email"`
	Code                 string
	Name                 string `json:"name"`
	Surname              string `json:"surname"`
	Patronymic           string `json:"patronymic"`
	RegistrationDateTime time.Time
	EncryptedPassword    string `json:"encrypted_password"`
	Role                 string `json:"role"`
	SupervisorID         string `json:"supervisor_id"`
	SupervisorInitials   string
}
