package server

import (
	"boom/handlers"
	"database/sql"

	"boom/proxy"

	"html/template"

	"github.com/gorilla/mux"
)

// RegisterRoutes регистрирует маршруты через Gorilla mux
func RegisterRoutes(m *mux.Router, tmpl *template.Template, dbConn *sql.DB) {

	// 🛰️ Подключаем все прокси
	proxy.ControlProxy(m)

	m.HandleFunc("/", handlers.BaseHandler(tmpl))

	m.HandleFunc("/theme", handlers.ToggleThemeHandler)

	m.HandleFunc("/spen", handlers.SpenHandler(tmpl, dbConn)).Methods("GET")
	m.HandleFunc("/spen", handlers.SpenAddHandler(tmpl, dbConn)).Methods("POST")
	m.HandleFunc("/delete", handlers.DeleteLastTodayHandler(tmpl, dbConn)).Methods("POST")

	m.HandleFunc("/deposit", handlers.DepositHandler(tmpl, dbConn)).Methods("GET")
	m.HandleFunc("/deposit", handlers.DepositAddHandler(tmpl, dbConn)).Methods("POST")

	m.HandleFunc("/delete-deposit", handlers.DeleteLastDepositTodayHandler(tmpl, dbConn)).Methods("POST")
}
