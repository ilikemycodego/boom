package server

import (
	"boom/db"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// NewServer собирает шаблоны, роуты и возвращает готовый http.Handler

func NewServer() http.Handler {
	tmpl := template.Must(
		template.New("").Funcs(template.FuncMap{}).ParseGlob("templates/**/*.html"),
	)

	dbConn := db.OpenDB("data.db")

	r := mux.NewRouter()
	r.Use(RequestLogger)

	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	RegisterRoutes(r, tmpl, dbConn)

	return r
}
