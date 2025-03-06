package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
)

const (
	driver   = "postgres"
	host     = "POSTGRES_DB_HOST"
	port     = "POSTGRES_DB_PORT"
	user     = "POSTGRES_DB_USER"
	password = "POSTGRES_DB_PASSWORD"
	dbname   = "POSTGRES_DB_NAME"
	sslmode  = "POSTGRES_DB_SSLMODE"
)

func NewPostgresConnection() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv(host),
		os.Getenv(port),
		os.Getenv(user),
		os.Getenv(password),
		os.Getenv(dbname),
		os.Getenv(sslmode))

	db, err := sql.Open(driver, dsn)
	if err != nil {
		slog.Error("open connection to postgres database failed")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		slog.Error("ping to postgres database failed")
		return nil, err
	}

	slog.Info("connected to postgres database",
		"host", os.Getenv(host),
		"port", os.Getenv(port),
	)
	return db, nil
}
