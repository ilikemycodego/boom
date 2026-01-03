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

// Данные для шаблона deposit (можно расширять)
type DepositData struct {
	TotalDay   int
	TotalWeek  int
	TotalMonth int
	TotalYear  int
}

// общий рендер
func renderDeposit(tmpl *template.Template, w http.ResponseWriter, data DepositData) {
	if err := tmpl.ExecuteTemplate(w, "deposit", data); err != nil {
		log.Println("❌ template deposit:", err)
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

// DepositHandler — GET: рендерит блок deposit
func DepositHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := loadDepositTotals(db)
		if err != nil {
			log.Println("❌ loadDepositTotals:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		renderDeposit(tmpl, w, data)
	}
}

// DepositAddHandler — POST: добавляет deposit и перерендерит блок deposit
func DepositAddHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad form", http.StatusBadRequest)
			return
		}

		amountStr := strings.TrimSpace(r.FormValue("amount"))
		amountStr = strings.ReplaceAll(amountStr, ",", ".")
		if strings.Contains(amountStr, ".") {
			http.Error(w, "только целое число", http.StatusBadRequest)
			return
		}

		amount, err := strconv.Atoi(amountStr)
		if err != nil || amount <= 0 {
			http.Error(w, "введите сумму > 0", http.StatusBadRequest)
			return
		}

		day := time.Now().Format("2006-01-02")

		if _, err := db.Exec(`INSERT INTO deposits(day, amount) VALUES(?, ?)`, day, amount); err != nil {
			log.Println("❌ insert deposit:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		data, err := loadDepositTotals(db)
		if err != nil {
			log.Println("❌ loadDepositTotals:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		renderDeposit(tmpl, w, data)
	}
}

// loadDepositTotals — суммы deposit день/неделя/месяц/год
func loadDepositTotals(db *sql.DB) (DepositData, error) {
	var out DepositData
	day := time.Now().Format("2006-01-02")

	// День
	if err := db.QueryRow(`SELECT COALESCE(SUM(amount),0) FROM deposits WHERE day = ?`, day).
		Scan(&out.TotalDay); err != nil {
		return out, err
	}

	// Неделя (календарная)
	if err := db.QueryRow(`
SELECT COALESCE(SUM(amount),0)
FROM deposits
WHERE strftime('%Y-%W', day) = strftime('%Y-%W', 'now', 'localtime')
`).Scan(&out.TotalWeek); err != nil {
		return out, err
	}

	// Месяц
	if err := db.QueryRow(`
SELECT COALESCE(SUM(amount),0)
FROM deposits
WHERE strftime('%Y-%m', day) = strftime('%Y-%m', 'now', 'localtime')
`).Scan(&out.TotalMonth); err != nil {
		return out, err
	}

	// Год
	if err := db.QueryRow(`
SELECT COALESCE(SUM(amount),0)
FROM deposits
WHERE strftime('%Y', day) = strftime('%Y', 'now', 'localtime')
`).Scan(&out.TotalYear); err != nil {
		return out, err
	}

	return out, nil
}
