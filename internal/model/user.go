package model

import "time"

type User struct {
	ID                   int       `json:"id"`
	Email                string    `json:"email"`
	Name                 string    `json:"name"`
	Surname              string    `json:"surname"`
	Patronymic           string    `json:"patronymic"`
	RegistrationDateTime time.Time `json:"reg_date_time"`
	EncryptedPassword    string    `json:"encrypted_password"`
	Role                 string    `json:"role"`
	SupervisorID         int       `json:"supervisor_id"`
}
