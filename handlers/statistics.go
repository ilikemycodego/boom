package handlers

import (
	"database/sql"
	"html/template"
	"log"

	"net/http"
)

// Данные для экрана статистики (и траты, и депозиты)
type StatisticsPageData struct {
	Spen    SpenData
	Deposit DepositData
}

func StatisticsHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		spen, err := loadTotals(db)
		if err != nil {
			log.Println("❌ loadTotals:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		dep, err := loadDepositTotals(db)
		if err != nil {
			log.Println("❌ loadDepositTotals:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		page := StatisticsPageData{
			Spen:    spen,
			Deposit: dep,
		}

		if err := tmpl.ExecuteTemplate(w, "statistics", page); err != nil {
			log.Printf("[StatisticsHandler] ❌ Ошибка шаблона: %v", err)
			http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
			return
		}
	}
}
