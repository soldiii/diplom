package repository

import (
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable       = "users"
	agentsTable      = "agents"
	supervisorsTable = "supervisors"
)

type PostgresDB struct {
	DatabaseURL string
	Password    string
	Database    *sqlx.DB
}

func NewPostgresDB(dbURL string) *PostgresDB {
	return &PostgresDB{
		DatabaseURL: dbURL,
		Password:    os.Getenv("DB_PASSWORD"),
	}
}

func (p *PostgresDB) OpenPostgresDB() (*sqlx.DB, error) {
	dbURL := p.DatabaseURL + " password=" + p.Password
	db, err := sqlx.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	p.Database = db
	return db, nil
}

func (p *PostgresDB) ClosePostgresDB() {
	p.Database.Close()
}
