package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/soldiii/diplom/internal/model"
)

type Authorization interface {
	CreateUser(user *model.User) (int, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
