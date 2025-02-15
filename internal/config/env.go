package config

import (
	"os"
	"strconv"
)

type Env struct {
	JWTSecret        string
	JWTExpireSeconds int
}

var Envs = initEnv()

func initEnv() Env {
	//godotenv.Load() //загружает данные из .env

	return Env{
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
