package repository

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresDB struct {
	DatabaseURL string
	Database    *sqlx.DB
}

func NewPostgresDB(dbURL string) *PostgresDB {
	return &PostgresDB{
		DatabaseURL: dbURL,
	}
}

func (p *PostgresDB) OpenPostgresDB() error {
	db, err := sqlx.Open("pgx", p.DatabaseURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	p.Database = db
	return nil
}

func (p *PostgresDB) ClosePostgresDB() {
	p.Database.Close()
}
