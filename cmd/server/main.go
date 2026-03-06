package main

import (
	"log"
	"net/http"

	"github.com/Kavibharathi-K/lifemetrics/internal/calories"
	"github.com/Kavibharathi-K/lifemetrics/internal/database"
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

	repo := calories.NewRepository(db)
	handler := calories.NewHandler(repo)

	r := chi.NewRouter()

	r.Post("/foods", handler.CreateFood)
	r.Get("/foods", handler.ListFoods)

	log.Println("🚀 Server running on :8080")

	http.ListenAndServe(":8080", r)
}
