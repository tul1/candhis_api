package db

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDatabaseConnection(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(2)
	db.SetConnMaxIdleTime(1 * time.Minute)
	db.SetConnMaxLifetime(3 * time.Minute)

	return db, nil
}
