package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBConnector func(dbURL string) (*sql.DB, error)

func DefaultDBConnector(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(2)
	db.SetConnMaxIdleTime(1 * time.Minute)
	db.SetConnMaxLifetime(3 * time.Minute)

	return db, nil
}

type DB struct {
	*sql.DB
	logger Logger
}

type Logger interface {
	Errorf(format string, args ...interface{})
}

func (db *DB) CloseWithLog() {
	if err := db.Close(); err != nil {
		db.logger.Errorf("Failed closing database: %v", err)
	}
}

func NewDBConnection(user, password, host, port, dbName string, connector DBConnector, logger Logger) (*DB, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)

	dbConn, err := connector(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return &DB{dbConn, logger}, nil
}
