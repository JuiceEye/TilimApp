package db

import (
	"database/sql"
	"fmt"
)

type DBConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}

func NewDBConnection(cfg DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		cfg.User, cfg.Password, cfg.Name, cfg.Host, cfg.Port)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to DB")
	return db, nil
}
