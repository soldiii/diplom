package model

import "time"

type User struct {
	Id                int       `json:"id"`
	Email             string    `json:"email"`
	Name              string    `json:"name"`
	Surname           string    `json:"surname"`
	Patronymic        string    `json:"patronymic"`
	Reg_date_time     time.Time `json:"reg_date_time"`
	EncryptedPassword string    `json:"encrypted_password"`
}
