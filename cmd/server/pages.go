package main

import (
	"html/template"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, page string, data interface{}) {

	tmpl := template.Must(template.ParseFiles(
		"assets/templates/layout.html",
		"assets/templates/"+page,
	))

	tmpl.ExecuteTemplate(w, "layout", data)
}

func HomePage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"ActivePage": "home",
		"SubPage":    "today",
	}

	renderTemplate(w, "home.html", data)
}

// Available Foods page
func FoodPage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"ActivePage": "food",
		"SubPage":    "foods",   // IMPORTANT
	}

	renderTemplate(w, "food.html", data)
}

// Meals page
func MealsPage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"ActivePage": "food",
		"SubPage":    "meals",
	}

	renderTemplate(w, "food.html", data)
}

func GoalsPage(
	w http.ResponseWriter,
	r *http.Request,
) {

	data := map[string]interface{}{

		"ActivePage": "home",

		"SubPage": "goals",
	}

	renderTemplate(

		w,

		"goals.html",

		data,
	)
}
