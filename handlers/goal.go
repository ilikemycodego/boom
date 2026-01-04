package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
)

type GoalData struct {
	Target  int
	Current int
	Percent int
}

func GoalHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			val, _ := strconv.Atoi(r.FormValue("target"))
			db.Exec("UPDATE goals SET amount = ? WHERE id = 1", val)
		}

		var target int
		db.QueryRow("SELECT amount FROM goals WHERE id = 1").Scan(&target)

		var current int
		// Считаем все депозиты за все время
		db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM deposits").Scan(&current)

		percent := 0
		if target > 0 {
			percent = (current * 100) / target
			if percent > 100 {
				percent = 100
			}
		}

		tmpl.ExecuteTemplate(w, "goal", GoalData{
			Target:  target,
			Current: current,
			Percent: percent,
		})
	}
}
