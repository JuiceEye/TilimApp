package main

import (
	"log"
	"tilimauth/cmd/api"
)

func main() {
	server := api.NewServer(":8080", nil)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

	// Database connection setup
	//connStr := "user=username dbname=mydb sslmode=disable"
	//db, err := sql.Open("postgres", connStr)
	//if err != nil {
	//	log.Fatal("Failed to connect to the database:", err)
	//}
	//defer db.Close()
	//
	//userRepo := &repository.UserRepository{DB: db}
	//userService := services.UserService{UserRepo: userRepo}
	//userHandler := handlers.UserHandler{UserService: userService}
	//
	//r := routes.NewRouter()
	//
	//log.Println("Server is running on :8080")
	//if err := http.ListenAndServe(":8080", routes); err != nil {
	//	log.Fatal("Failed to start server:", err)
	//}
}
