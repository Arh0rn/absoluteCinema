package database

import (
	"database/sql"
	"fmt"
	"os"
)

type Config struct {
	Host    string
	Port    string
	User    string
	DBName  string
	SSLMode string
}

func NewPostgresConnection(config Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, os.Getenv("postgres"), config.DBName, config.SSLMode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
