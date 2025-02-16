package main

import (
	"log"
	"tilimauth/cmd/api"
	"tilimauth/internal/config"
	"tilimauth/internal/db"
)

func main() {
	//todo: пофиксить гавно сделанное для логгирования кек
	config1 := db.Config{
		User:     config.Envs.DBUser,
		Password: config.Envs.DBPassword,
		Name:     config.Envs.DBName,
		Host:     config.Envs.DBHost,
		Port:     config.Envs.DBPort,
	}
	dbase, err := db.NewDBConnection(config1)
	log.Printf("Connecting to db %v", config1)

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(":8080", dbase)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
