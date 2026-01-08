package server

import (
	"boom/food"

	"database/sql"

	"html/template"

	"github.com/gorilla/mux"
)

// RegisterRoutes регистрирует маршруты через Gorilla mux
func RoutesFood(m *mux.Router, tmpl *template.Template, dbConn *sql.DB) {

	m.HandleFunc("/food", food.FoodHandler(tmpl))

}
