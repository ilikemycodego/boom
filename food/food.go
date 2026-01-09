package food

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// FoodHandler рендерит страницу "Еда" (облако меток + сохранение)
func FoodHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		day := todayStr()

		tags, err := ListTags(db)
		if err != nil {
			log.Printf("[FoodHandler] ❌ ListTags: %v", err)
			http.Error(w, "Ошибка загрузки меток", http.StatusInternalServerError)
			return
		}

		selected, err := GetSelectedTagIDs(db, day)
		if err != nil {
			log.Printf("[FoodHandler] ❌ GetSelectedTagIDs: %v", err)
			http.Error(w, "Ошибка загрузки выбранных меток", http.StatusInternalServerError)
			return
		}

		for i := range tags {
			tags[i].Selected = selected[tags[i].ID]
		}

		data := FoodPageData{
			Today: day,
			Tags:  tags,
		}

		if err := tmpl.ExecuteTemplate(w, "food", data); err != nil {
			log.Printf("[FoodHandler] ❌ template food: %v", err)
			http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
			return
		}
	}
}

// FoodSaveEntryHandler сохраняет выбор меток за сегодня
func FoodSaveEntryHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка формы", http.StatusBadRequest)
			return
		}

		day := strings.TrimSpace(r.FormValue("entry_date"))
		if day == "" {
			day = todayStr()
		}

		// Собираем id меток
		var tagIDs []int
		for _, v := range r.Form["tag_id"] {
			id, err := strconv.Atoi(v)
			if err != nil || id <= 0 {
				continue
			}
			tagIDs = append(tagIDs, id)
		}

		if err := SaveEntry(db, day, tagIDs); err != nil {
			log.Printf("[FoodSaveEntryHandler] ❌ SaveEntry: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			_ = tmpl.ExecuteTemplate(w, "food_save_result", map[string]any{
				"OK":      false,
				"Message": "Не удалось сохранить",
			})
			return
		}

		_ = tmpl.ExecuteTemplate(w, "food_save_result", map[string]any{
			"OK":      true,
			"Message": "Сохранено ✅",
		})
	}
}

// FoodTagsHandler:
// - GET  /food/tags  -> вернуть панель управления метками
// - POST /food/tags  -> добавить метку + обновить панель (и облако через OOB)
func FoodTagsHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Добавление
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err == nil {
				name := r.FormValue("name")
				if err := AddTag(db, name); err != nil {
					// Можно вернуть 400, если хочешь показывать ошибку пользователю
					log.Printf("[FoodTagsHandler] ❌ AddTag: %v", err)
				}
			}
		}

		tags, err := ListTags(db)
		if err != nil {
			http.Error(w, "Ошибка загрузки меток", http.StatusInternalServerError)
			return
		}

		// Панель управления всегда возвращаем в тело ответа
		if err := tmpl.ExecuteTemplate(w, "food_tags_panel", map[string]any{
			"Tags": tags,
		}); err != nil {
			log.Printf("[FoodTagsHandler] ❌ template food_tags_panel: %v", err)
		}

		// Для POST дополнительно обновим облако меток (hx-swap-oob)
		if r.Method == http.MethodPost {
			_ = tmpl.ExecuteTemplate(w, "food_cloud_oob", FoodPageData{
				Today: todayStr(),
				Tags:  tags,
			})
		}
	}
}

// FoodTagDeleteHandler удаляет метку + обновляет панель и облако (OOB)
func FoodTagDeleteHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil || id <= 0 {
			http.Error(w, "Некорректный id", http.StatusBadRequest)
			return
		}

		if err := DeleteTag(db, id); err != nil {
			log.Printf("[FoodTagDeleteHandler] ❌ DeleteTag: %v", err)
		}

		tags, err := ListTags(db)
		if err != nil {
			http.Error(w, "Ошибка загрузки меток", http.StatusInternalServerError)
			return
		}

		// Обновляем панель
		if err := tmpl.ExecuteTemplate(w, "food_tags_panel", map[string]any{
			"Tags": tags,
		}); err != nil {
			log.Printf("[FoodTagDeleteHandler] ❌ template food_tags_panel: %v", err)
		}

		// И облако (OOB)
		_ = tmpl.ExecuteTemplate(w, "food_cloud_oob", FoodPageData{
			Today: todayStr(),
			Tags:  tags,
		})
	}
}
