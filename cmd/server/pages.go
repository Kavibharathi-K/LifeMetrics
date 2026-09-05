package main

import (
	"html/template"
	"net/http"
)

func renderTemplate(
	w http.ResponseWriter,
	page string,
	data interface{},
) {

	tmpl := template.Must(

		template.ParseFiles(

			"assets/templates/layout.html",

			"assets/templates/"+page,

			"assets/templates/partials/navbar.html",

			"assets/templates/partials/submenu.html",

			"assets/templates/partials/scripts.html",
		),
	)

	tmpl.ExecuteTemplate(
		w,
		"layout",
		data,
	)
}

func HomePage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"ActivePage": "home",
		"SubPage":    "today",
	}

	renderTemplate(w, "home.html", data)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(
		template.ParseFiles(
			"assets/templates/login.html",
		),
	)

	tmpl.ExecuteTemplate(
		w,
		"login",
		nil,
	)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(
		template.ParseFiles(
			"assets/templates/register.html",
		),
	)

	tmpl.ExecuteTemplate(
		w,
		"register",
		nil,
	)
}


// Available Foods page
func FoodPage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"ActivePage": "food",
		"SubPage":    "foods",
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

func WorkoutPage(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"ActivePage": "workout",
		"SubPage":    "workout",
	}

	renderTemplate(w, "workout.html", data)
}
