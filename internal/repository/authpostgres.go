package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/soldiii/diplom/internal/model"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user *model.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (email, name, surname, patronymic, reg_date_time, encrypted_password, role) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	row := r.db.QueryRow(query, user.Email, user.Name, user.Surname, user.Patronymic, user.RegistrationDateTime, user.EncryptedPassword, user.Role)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
