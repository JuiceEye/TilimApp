package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Env struct {
	DBUser                  string
	DBPassword              string
	DBName                  string
	DBHost                  string
	DBPort                  string
	JWTSecret               string
	JWTAccessExpireSeconds  int
	JWTRefreshExpireSeconds int
}

var Envs = initEnv()

func initEnv() Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("[ERROR] Could not load .env file: %v", err)
	}

	return Env{
		DBUser:                  getEnv("DB_USER", ""),
		DBPassword:              getEnv("DB_PASSWORD", ""),
		DBName:                  getEnv("DB_NAME", ""),
		DBHost:                  getEnv("DB_HOST", ""),
		DBPort:                  getEnv("DB_PORT", ""),
		JWTSecret:               getEnv("JWT_SECRET", ""),
		JWTAccessExpireSeconds:  getEnvAsInt("JWT_ACCESS_EXPIRE_SECONDS", 300),         // 5 минут аксесс токен
		JWTRefreshExpireSeconds: getEnvAsInt("JWT_REFRESH_EXPIRE_SECONDS", 3600*24*30), // 30 дней рефреш токен
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
