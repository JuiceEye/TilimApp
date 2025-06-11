package db

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
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
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
		cfg.user, cfg.password, cfg.host, cfg.port, cfg.name)

	pgxCfg, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	pgxCfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	db := stdlib.OpenDB(*pgxCfg)

	log.Printf(`[INFO] Connected to db "%s" %s:%s`, cfg.name, cfg.host, cfg.port)
	return db, nil
}
