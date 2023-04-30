package model

import "time"

type User struct {
	ID                   int
	Email                string
	Name                 string
	Surname              string
	Patronymic           string
	RegistrationDateTime time.Time
	EncryptedPassword    string
	Role                 string
	SupervisorID         string
}
