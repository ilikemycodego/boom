package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// данные для шаблона spen
type SpenData struct {
	TotalDay   int
	TotalWeek  int
	TotalMonth int
	TotalYear  int
}

// SpenHandler — рендерит блок с суммами
func SpenHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := loadTotals(db)
		if err != nil {
			log.Println("❌ loadTotals:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "spen", data); err != nil {
			log.Printf("[SpenHandler] ❌ Ошибка шаблона: %v", err)
			http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
			return
		}
	}
}

// SpenAddHandler — добавляет трату

func SpenAddHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form", http.StatusBadRequest)
			return
		}

		amountStr := strings.TrimSpace(r.FormValue("amount"))
		amountStr = strings.ReplaceAll(amountStr, ",", ".") // если вдруг введут 10,5

		// Самый простой вариант: принимаем только целые
		if strings.Contains(amountStr, ".") {
			http.Error(w, "только целое число (без .00)", http.StatusBadRequest)
			return
		}

		amount, err := strconv.Atoi(amountStr)
		if err != nil || amount <= 0 {
			http.Error(w, "введите сумму > 0", http.StatusBadRequest)
			return
		}

		day := time.Now().Format("2006-01-02")

		if _, err := db.Exec(`INSERT INTO expenses(day, amount) VALUES(?, ?)`, day, amount); err != nil {
			log.Println("❌ insert:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		// после вставки — перерендерим тот же блок spen
		data, err := loadTotals(db)
		if err != nil {
			log.Println("❌ loadTotals:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "spen", data); err != nil {
			log.Println("❌ template:", err)
			http.Error(w, "template error", http.StatusInternalServerError)
			return
		}
	}
}

// loadTotals считает суммы на основе сохранённых расходов
func loadTotals(db *sql.DB) (SpenData, error) {
	var out SpenData
	day := time.Now().Format("2006-01-02")

	// День
	if err := db.QueryRow(`SELECT COALESCE(SUM(amount),0) FROM expenses WHERE day = ?`, day).Scan(&out.TotalDay); err != nil {
		return out, err
	}

	// Неделя/месяц/год — считаем по day (TEXT YYYY-MM-DD)
	// SQLite: week = strftime('%Y-%W', day)
	if err := db.QueryRow(`
SELECT COALESCE(SUM(amount),0)
FROM expenses
WHERE strftime('%Y-%W', day) = strftime('%Y-%W', 'now', 'localtime')
`).Scan(&out.TotalWeek); err != nil {
		return out, err
	}

	if err := db.QueryRow(`
SELECT COALESCE(SUM(amount),0)
FROM expenses
WHERE strftime('%Y-%m', day) = strftime('%Y-%m', 'now', 'localtime')
`).Scan(&out.TotalMonth); err != nil {
		return out, err
	}

	if err := db.QueryRow(`
SELECT COALESCE(SUM(amount),0)
FROM expenses
WHERE strftime('%Y', day) = strftime('%Y', 'now', 'localtime')
`).Scan(&out.TotalYear); err != nil {
		return out, err
	}

	return out, nil
}
