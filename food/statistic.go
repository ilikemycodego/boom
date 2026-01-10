package food

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// FoodDayStat — строка статистики по одному дню
type FoodDayStat struct {
	Day  string
	Tags string // строка "хлеб, мясо, гречка"
}

// ListFoodStats возвращает список дней и меток (сгруппировано по дню)
func ListFoodStats(db *sql.DB) ([]FoodDayStat, error) {
	rows, err := db.Query(`
SELECT
  e.entry_date,
  COALESCE(GROUP_CONCAT(t.name, ', '), '') AS tags
FROM food_entries e
LEFT JOIN food_entry_tags et ON et.entry_id = e.id
LEFT JOIN food_tags t ON t.id = et.tag_id
GROUP BY e.id
ORDER BY e.id DESC;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []FoodDayStat
	for rows.Next() {
		var s FoodDayStat
		if err := rows.Scan(&s.Day, &s.Tags); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// StatisticHandler рендерит страницу статистики по еде
func StatisticHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := ListFoodStats(db)
		if err != nil {
			log.Printf("[StatisticHandler] ❌ ListFoodStats: %v", err)
			http.Error(w, "Ошибка загрузки статистики", http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"Stats": stats,
		}

		if err := tmpl.ExecuteTemplate(w, "food__statistics", data); err != nil {
			log.Printf("[StatisticHandler] ❌ template food__statistics: %v", err)
			http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
			return
		}
	}
}
