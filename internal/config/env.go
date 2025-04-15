package config

import (
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
	//godotenv.Load() 	//загружает данные из .env

	//todo: починить доставание из env
	return Env{
		DBUser:                  getEnv("DB_USER", "postgres.nffhzpzmtuicxmipwlcz"),
		DBPassword:              getEnv("DB_PASSWORD", "farukhnastya2003"),
		DBName:                  getEnv("DB_NAME", "auth"),
		DBHost:                  getEnv("DB_HOST", "aws-0-eu-central-1.pooler.supabase.com"),
		DBPort:                  getEnv("DB_PORT", "6543"),
		JWTSecret:               getEnv("JWT_SECRET", "supersecretkey"),
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
