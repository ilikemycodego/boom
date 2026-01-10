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
	m.HandleFunc("/food/tags", food.FoodTagsPageHandler(tmpl, dbConn)).Methods("GET", "POST")
	m.HandleFunc("/food/tags/{id:[0-9]+}", food.FoodTagDeleteHandler(tmpl, dbConn)).Methods("DELETE")

	m.HandleFunc("/food/select", food.FoodSelectHandler(tmpl, dbConn)).Methods("POST")
	m.HandleFunc("/food/clear", food.FoodClearHandler(tmpl)).Methods("POST")

	m.HandleFunc("/food/statistics", food.StatisticHandler(tmpl, dbConn)).Methods("GET")
}
