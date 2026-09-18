package database

import (
	"database/sql"
	"fmt"
	_"github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	User 	string
	Password 	string
	Host 	string
	Port 	string
	DBname 	string
}

func Open(cfg Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBname,
	)

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}