package food

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

type FoodEntryRow struct {
	ID   int
	Time string
	Tags string
}

type FoodDayGroup struct {
	Day     string
	Entries []FoodEntryRow
}

type FoodStatisticsPageData struct {
	Days []FoodDayGroup
}

func ListFoodStatsGrouped(db *sql.DB) ([]FoodDayGroup, error) {
	rows, err := db.Query(`
SELECT
  e.entry_date,
  e.id,
  strftime('%H:%M', e.created_at) AS created_time,
  COALESCE(
    (
      SELECT GROUP_CONCAT(x.name, ', ')
      FROM (
        SELECT t2.name AS name
        FROM food_entry_tags et2
        JOIN food_tags t2 ON t2.id = et2.tag_id
        WHERE et2.entry_id = e.id
        ORDER BY t2.name
      ) AS x
    ),
    ''
  ) AS tags
FROM food_entries e
ORDER BY e.entry_date DESC, e.id DESC;
`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []FoodDayGroup
	var cur *FoodDayGroup

	for rows.Next() {
		var day string
		var id int
		var tm string
		var tags string

		if err := rows.Scan(&day, &id, &tm, &tags); err != nil {
			return nil, err
		}

		if cur == nil || cur.Day != day {
			out = append(out, FoodDayGroup{Day: day})
			cur = &out[len(out)-1]
		}

		cur.Entries = append(cur.Entries, FoodEntryRow{
			ID:   id,
			Time: tm,
			Tags: tags,
		})
	}

	return out, rows.Err()
}
func StatisticHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		days, err := ListFoodStatsGrouped(db)
		if err != nil {
			log.Printf("[StatisticHandler] ❌ ListFoodStatsGrouped: %v", err)
			http.Error(w, "Ошибка загрузки статистики", http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "food__statistics", FoodStatisticsPageData{
			Days: days,
		}); err != nil {
			log.Printf("[StatisticHandler] ❌ template food__statistics: %v", err)
			http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
			return
		}
	}
}
