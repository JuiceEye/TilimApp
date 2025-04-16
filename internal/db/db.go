package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"tilimauth/internal/config"
)

type Config struct {
	user     string
	password string
	name     string
	host     string
	port     string
}

func NewConfig(envs *config.Env) *Config {
	return &Config{
		user:     envs.DBUser,
		password: envs.DBPassword,
		name:     envs.DBName,
		host:     envs.DBHost,
		port:     envs.DBPort,
	}
}

func NewDBConnection(cfg *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=require",
		cfg.user, cfg.password, cfg.name, cfg.host, cfg.port)
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	log.Printf(`[INFO] Connected to db "%s" %s:%s`, cfg.name, cfg.host, cfg.port)

	return db, nil
}
