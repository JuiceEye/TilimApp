package main

import (
	"log"
	"tilimauth/cmd/api"
	"tilimauth/internal/config"
	"tilimauth/internal/db"
)

func main() {
	dbConfig := db.NewConfig(&config.Envs)
	dbase, err := db.NewDBConnection(dbConfig)

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(":8080", dbase)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
