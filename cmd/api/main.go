package main

import (
	"REST-service-sub/internal/config"
	"REST-service-sub/internal/db"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	_ = godotenv.Load("../../.env")
	cfg := config.LoadConfig()
	fmt.Println("Config loaded successfully!")

	fmt.Println("DSN:", cfg.DSN())

	database, err := db.NewPostgresDB(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed connect to Postgres: %v", err)
	}
	defer database.Close()
	fmt.Println("Database connected successfully!")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	fmt.Printf("Server is listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
