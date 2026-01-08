package server

import (
	"boom/ui"
	"database/sql"

	"html/template"

	"github.com/gorilla/mux"
)

// RegisterRoutes регистрирует маршруты через Gorilla mux
func RoutesUI(m *mux.Router, tmpl *template.Template, dbConn *sql.DB) {

	m.HandleFunc("/", ui.BaseHandler(tmpl))

	m.HandleFunc("/theme", ui.ToggleThemeHandler)

	m.HandleFunc("/spen", ui.SpenHandler(tmpl, dbConn)).Methods("GET")
	m.HandleFunc("/spen", ui.SpenAddHandler(tmpl, dbConn)).Methods("POST")
	m.HandleFunc("/delete", ui.DeleteLastTodayHandler(tmpl, dbConn)).Methods("POST")

	m.HandleFunc("/deposit", ui.DepositHandler(tmpl, dbConn)).Methods("GET")
	m.HandleFunc("/deposit", ui.DepositAddHandler(tmpl, dbConn)).Methods("POST")

	m.HandleFunc("/delete-deposit", ui.DeleteLastDepositTodayHandler(tmpl, dbConn)).Methods("POST")

	m.HandleFunc("/goal", ui.GoalBarHandler(tmpl, dbConn)).Methods("GET")
	m.HandleFunc("/goal-target", ui.GoalTargetHandler(tmpl, dbConn)).Methods("GET", "POST")

	m.HandleFunc("/clear-all-data", ui.ClearAllDataHandler(dbConn)).Methods("POST")

	m.HandleFunc("/statistics", ui.StatisticsHandler(tmpl, dbConn)).Methods("GET")

}
