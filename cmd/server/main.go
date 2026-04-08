package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Kavibharathi-K/lifemetrics/database"
	"github.com/Kavibharathi-K/lifemetrics/services/food"
	"github.com/Kavibharathi-K/lifemetrics/services/meal"
	"github.com/Kavibharathi-K/lifemetrics/services/meal_item"
	"github.com/Kavibharathi-K/lifemetrics/services/user_metrics"
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
	meal_item.RegisterRoutes(r, db)
	meal.RegisterRoutes(r, db)
	user_metrics.RegisterRoutes(r, db)

	log.Println("=======================")
	log.Println("Life Metrics is running")
	log.Println("=======================")

	// pages
	r.Get("/", HomePage)
	r.Get("/food", FoodPage)
	r.Get("/meals", MealsPage)

	// static files
	fs := http.FileServer(http.Dir("./assets/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	wd, _ := os.Getwd()
	fmt.Println("Working dir:", wd)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
