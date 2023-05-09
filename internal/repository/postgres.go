package repository

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable       = "users"
	agentsTable      = "agents"
	supervisorsTable = "supervisors"
	userCodesTable   = "usercodes"
	adsTable         = "ads"
	plansTable       = "plans"
	reportsTable     = "reports"
)

type DatabaseURL struct {
	Host     string
	Port     string
	DBname   string
	Username string
	Password string
	SSLMode  string
}

func NewDatabaseURL() *DatabaseURL {
	return &DatabaseURL{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBname:   os.Getenv("DB_NAME"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
}

type PostgresDB struct {
	databaseURL *DatabaseURL
	Database    *sqlx.DB
}

func NewPostgresDB(dbURL *DatabaseURL) *PostgresDB {
	return &PostgresDB{
		databaseURL: dbURL,
	}
}

func RunMigration(dbURL string) error {
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) OpenPostgresDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		p.databaseURL.Host, p.databaseURL.Port, p.databaseURL.DBname, p.databaseURL.Username, p.databaseURL.Password, p.databaseURL.SSLMode))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	databaseURL := "postgres://" + p.databaseURL.Username + ":" + p.databaseURL.Password + "@" + p.databaseURL.Host + ":" + p.databaseURL.Port + "/" + p.databaseURL.DBname + "?sslmode=" + p.databaseURL.SSLMode
	RunMigration(databaseURL)
	p.Database = db
	return db, nil
}

func (p *PostgresDB) ClosePostgresDB() {
	p.Database.Close()
}
