package main

import (
	"log"
	"tilimauth/cmd/api"
	"tilimauth/internal/config"
	"tilimauth/internal/db"
)

func main() {
	dbase, err := db.NewDBConnection(db.DBConfig{
		User:     config.Envs.DBUser,
		Password: config.Envs.DBPassword,
		Name:     config.Envs.DBName,
		Host:     config.Envs.DBHost,
		Port:     config.Envs.DBUser,
	})

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(":8080", dbase)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
