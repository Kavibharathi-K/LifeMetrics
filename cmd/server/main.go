package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	log.Println("=======================")
	log.Println("Life Metrics is running")
	log.Println("=======================")

	fs := http.FileServer(http.Dir("./frontend"))

	r.Handle("/*", fs)

	wd, _ := os.Getwd()
	fmt.Println("Working dir:", wd)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
