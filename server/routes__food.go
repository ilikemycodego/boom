package server

import (
	"boom/food"

	"database/sql"

	"html/template"

	"github.com/gorilla/mux"
)

// RegisterRoutes регистрирует маршруты через Gorilla mux
func RoutesFood(m *mux.Router, tmpl *template.Template, dbConn *sql.DB) {

	m.HandleFunc("/food", food.FoodHandler(tmpl, dbConn)).Methods("GET")
	m.HandleFunc("/food/entry", food.FoodSaveEntryHandler(tmpl, dbConn)).Methods("POST")

	// Панель управления метками
	m.HandleFunc("/food/tags", food.FoodTagsHandler(tmpl, dbConn)).Methods("GET", "POST")
	m.HandleFunc("/food/tags/{id:[0-9]+}", food.FoodTagDeleteHandler(tmpl, dbConn)).Methods("DELETE")

}
