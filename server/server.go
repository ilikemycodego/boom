package server

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// Теперь функция принимает уже открытое соединение dbConn
func NewServer(dbConn *sql.DB) http.Handler {
	tmpl := template.Must(
		template.New("").Funcs(template.FuncMap{}).ParseGlob("templates/**/*.html"),
	)

	// -
	r := mux.NewRouter()
	r.Use(RequestLogger)

	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Передаем полученный dbConn дальше в роуты
	RoutesUI(r, tmpl, dbConn)
	RoutesFood(r, tmpl, dbConn)

	return r
}
