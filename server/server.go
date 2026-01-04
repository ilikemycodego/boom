package server

import (
	"database/sql" // Добавьте этот импорт
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// Теперь функция принимает уже открытое соединение dbConn
func NewServer(dbConn *sql.DB) http.Handler {
	tmpl := template.Must(
		template.New("").Funcs(template.FuncMap{}).ParseGlob("templates/**/*.html"),
	)

	// УДАЛИЛИ строки с db.OpenDB("data.db"), так как база приходит снаружи

	r := mux.NewRouter()
	r.Use(RequestLogger)

	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Передаем полученный dbConn дальше в роуты
	RegisterRoutes(r, tmpl, dbConn)

	return r
}
