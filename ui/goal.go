package ui

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

func loadGoal(db *sql.DB) (GoalData, error) {
	var out GoalData

	// цель
	if err := db.QueryRow("SELECT amount FROM goals WHERE id = 1").Scan(&out.Target); err != nil {
		return out, err
	}

	// текущая сумма (все депозиты)
	if err := db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM deposits").Scan(&out.Current); err != nil {
		return out, err
	}

	// процент
	out.Percent = 0
	if out.Target > 0 {
		out.Percent = (out.Current * 100) / out.Target
		if out.Percent > 100 {
			out.Percent = 100
		}
	}

	return out, nil
}

func GoalBarHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data, err := loadGoal(db)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "goal_bar", data)
	}
}

func GoalTargetHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			val, _ := strconv.Atoi(r.FormValue("target"))
			if _, err := db.Exec("UPDATE goals SET amount = ? WHERE id = 1", val); err != nil {
				http.Error(w, "DB error", http.StatusInternalServerError)
				return
			}

			// Сообщаем всем блокам (в т.ч. goal_bar), что цель обновилась
			w.Header().Set("HX-Trigger", "goalUpdated")
		}

		// Рендерим форму с текущей целью (удобно видеть текущее значение)
		data, err := loadGoal(db)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "goal_target", data)
	}
}
