package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"
)

// DeleteLastTodayHandler — удаляет последнюю трату за сегодня
func DeleteLastTodayHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		day := time.Now().Format("2006-01-02")

		_, err := db.Exec(`
DELETE FROM expenses
WHERE id = (
  SELECT id
  FROM expenses
  WHERE day = ?
  ORDER BY id DESC
  LIMIT 1
)`, day)
		if err != nil {
			log.Println("❌ delete last:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		// После удаления — на главную
		w.Header().Set("HX-Redirect", "/")
	}
}

// DeleteLastDepositTodayHandler — удаляет последний deposit за сегодня
func DeleteLastDepositTodayHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		day := time.Now().Format("2006-01-02")

		_, err := db.Exec(`
DELETE FROM deposits
WHERE id = (
  SELECT id
  FROM deposits
  WHERE day = ?
  ORDER BY id DESC
  LIMIT 1
)`, day)
		if err != nil {
			log.Println("❌ delete last deposit:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/")
	}
}

// удаляет все данные
func ClearAllDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Используем транзакцию, чтобы убедиться, что обе операции либо выполнятся, либо нет
		tx, err := db.Begin()
		if err != nil {
			log.Println("❌ clear all data begin transaction:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		// Очищаем таблицу расходов
		if _, err := tx.Exec(`DELETE FROM expenses`); err != nil {
			tx.Rollback() // Откатываем изменения, если что-то пошло не так
			log.Println("❌ clear all data delete expenses:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		// Очищаем таблицу депозитов
		if _, err := tx.Exec(`DELETE FROM deposits`); err != nil {
			tx.Rollback() // Откатываем изменения
			log.Println("❌ clear all data delete deposits:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		// Подтверждаем транзакцию
		if err := tx.Commit(); err != nil {
			log.Println("❌ clear all data commit:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		// После успешной очистки перенаправляем пользователя на главную страницу
		w.Header().Set("HX-Redirect", "/")
	}
}
