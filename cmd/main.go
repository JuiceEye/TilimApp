package main

import (
	"log"
	"tilimauth/cmd/api"
	"tilimauth/infrastructure"
	"tilimauth/internal/config"
	"tilimauth/internal/db"
)

func main() {
	dbConfig := db.NewConfig(&config.Envs)
	dbase, err := db.NewDBConnection(dbConfig)

	if err != nil {
		log.Fatal(err)
	}

	redis := infrastructure.NewRedisClient()

	server := api.NewServer("127.0.0.1:8080", dbase, redis)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
