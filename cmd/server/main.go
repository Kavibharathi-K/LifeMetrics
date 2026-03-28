package main

import (
	"log"
	"net/http"

	"github.com/Kavibharathi-K/lifemetrics/database"
	"github.com/Kavibharathi-K/lifemetrics/services/food"
	"github.com/go-chi/chi/v5"
)

func main() {

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}

	err = database.RunMigrations(db)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	food.RegisterRoutes(r, db)
	
	log.Println("🚀 Server running on :8080")
	http.ListenAndServe(":8080", r)
}
