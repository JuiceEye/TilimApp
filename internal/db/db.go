package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Config struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}

func NewDBConnection(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Name, cfg.Host, cfg.Port)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to DB")
	return db, nil
}
