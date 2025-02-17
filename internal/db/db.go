package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"tilimauth/internal/config"
)

type Config struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}

func NewConfig(envs *config.Env) *Config {
	return &Config{
		User:     envs.DBUser,
		Password: envs.DBPassword,
		Name:     envs.DBName,
		Host:     envs.DBHost,
		Port:     envs.DBPort,
	}
}

func NewDBConnection(cfg *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Name, cfg.Host, cfg.Port)
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	log.Printf("Connected to db %v", cfg)

	return db, nil
}
