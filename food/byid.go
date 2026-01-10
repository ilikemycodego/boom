package food

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// TagsByIDs возвращает теги по списку id (для превью "выбрано")
func TagsByIDs(db *sql.DB, ids []int) ([]Tag, error) {
	if len(ids) == 0 {
		return []Tag{}, nil
	}

	// (?, ?, ?, ...) под количество ids
	ph := strings.TrimRight(strings.Repeat("?,", len(ids)), ",")
	q := fmt.Sprintf(`SELECT id, name FROM food_tags WHERE id IN (%s) ORDER BY name ASC`, ph)

	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}

	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Tag, 0, len(ids))
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// FoodSelectHandler: добавить тег в список (без toggle)
func FoodSelectHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка формы", http.StatusBadRequest)
			return
		}

		// Текущие выбранные (из hidden input name="tag_id")
		var ids []int
		for _, v := range r.Form["tag_id"] {
			id, err := strconv.Atoi(v)
			if err == nil && id > 0 {
				ids = append(ids, id)
			}
		}

		// Добавляем новый (из hx-vals)
		addID, _ := strconv.Atoi(r.FormValue("add_tag_id"))
		if addID > 0 {
			ids = append(ids, addID)
		}

		// Убираем дубликаты
		seen := map[int]bool{}
		uniq := make([]int, 0, len(ids))
		for _, id := range ids {
			if id <= 0 || seen[id] {
				continue
			}
			seen[id] = true
			uniq = append(uniq, id)
		}

		// Чтоб результат был стабильный (не обязательно)
		sort.Ints(uniq)

		tags, err := TagsByIDs(db, uniq)
		if err != nil {
			log.Printf("[FoodSelectHandler] ❌ TagsByIDs: %v", err)
			http.Error(w, "Ошибка выбора", http.StatusInternalServerError)
			return
		}

		_ = tmpl.ExecuteTemplate(w, "food_selected", tags)
	}
}

// FoodClearHandler: полностью очистить список (и hidden inputs)
func FoodClearHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = tmpl.ExecuteTemplate(w, "food_selected", []Tag{})
	}
}
