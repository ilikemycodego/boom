package server

import (
	"boom/handlers"
	"database/sql"

	"html/template"

	"github.com/gorilla/mux"
)

// RegisterRoutes регистрирует маршруты через Gorilla mux
func RegisterRoutes(m *mux.Router, tmpl *template.Template, dbConn *sql.DB) {

	m.HandleFunc("/", handlers.BaseHandler(tmpl))

	m.HandleFunc("/theme", handlers.ToggleThemeHandler)

	m.HandleFunc("/spen", handlers.SpenHandler(tmpl, dbConn)).Methods("GET")
	m.HandleFunc("/spen", handlers.SpenAddHandler(tmpl, dbConn)).Methods("POST")
	m.HandleFunc("/delete", handlers.DeleteLastTodayHandler(tmpl, dbConn)).Methods("POST")

	m.HandleFunc("/deposit", handlers.DepositHandler(tmpl, dbConn)).Methods("GET")
	m.HandleFunc("/deposit", handlers.DepositAddHandler(tmpl, dbConn)).Methods("POST")

	m.HandleFunc("/delete-deposit", handlers.DeleteLastDepositTodayHandler(tmpl, dbConn)).Methods("POST")

	m.HandleFunc("/goal", handlers.GoalBarHandler(tmpl, dbConn)).Methods("GET")
	m.HandleFunc("/goal-target", handlers.GoalTargetHandler(tmpl, dbConn)).Methods("GET", "POST")

	m.HandleFunc("/clear-all-data", handlers.ClearAllDataHandler(dbConn)).Methods("POST")

	m.HandleFunc("/statistics", handlers.StatisticsHandler(tmpl, dbConn)).Methods("GET")

}
