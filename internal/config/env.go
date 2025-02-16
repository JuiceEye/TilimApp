package config

import (
	"os"
	"strconv"
)

type Env struct {
	DBUser           string
	DBPassword       string
	DBName           string
	DBHost           string
	DBPort           string
	JWTSecret        string
	JWTExpireSeconds int
}

var Envs = initEnv()

func initEnv() Env {
	//godotenv.Load() 	//загружает данные из .env

	return Env{
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTExpireSeconds: getEnvAsInt("JWT_EXPIRE_SECONDS", 3600*24*7), //7 дней
	}
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
