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
			http.Error(w, "Ошибка загрузки меток", http.StatusInternalServerError)
			return
		}

		data := FoodPageData{
			Today:        day,
			Tags:         tags,
			SelectedTags: []Tag{}, // ✅ всегда пусто
		}

		if err := tmpl.ExecuteTemplate(w, "food", data); err != nil {
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
			http.Error(w, "Не удалось сохранить", http.StatusBadRequest)
			return
		}

		// ✅ После сохранения: очищаем "форму" (выбор) — рендерим food с пустым selected
		tags, err := ListTags(db)
		if err != nil {
			log.Printf("[FoodSaveEntryHandler] ❌ ListTags: %v", err)
			http.Error(w, "Ошибка загрузки меток", http.StatusInternalServerError)
			return
		}

		if len(tagIDs) == 0 {
			// Ничего не сохраняем, просто перерисуем страницу пустой
			tags, _ := ListTags(db)
			_ = tmpl.ExecuteTemplate(w, "food", FoodPageData{
				Today:        day,
				Tags:         tags,
				SelectedTags: []Tag{},
			})
			return
		}

		data := FoodPageData{
			Today:        day,
			Tags:         tags,
			SelectedTags: []Tag{}, // ✅ пусто = UI чистый
		}

		if err := tmpl.ExecuteTemplate(w, "food", data); err != nil {
			log.Printf("[FoodSaveEntryHandler] ❌ template food: %v", err)
			http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
			return
		}
	}
}

// FoodTagsHandler:
// - GET  /food/tags  -> вернуть панель управления метками
// - POST /food/tags  -> добавить метку + обновить панель (и облако через OOB)
// FoodTagsPageHandler:
// - GET  /food/tags  -> показать страницу управления
// - POST /food/tags  -> добавить метку и перерисовать страницу
func FoodTagsPageHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Добавление
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err == nil {
				name := r.FormValue("name")
				if err := AddTag(db, name); err != nil {
					log.Printf("[FoodTagsPageHandler] ❌ AddTag: %v", err)
				}
			}
		}

		tags, err := ListTags(db)
		if err != nil {
			http.Error(w, "Ошибка загрузки меток", http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "food_tags_page", map[string]any{
			"Tags": tags,
		}); err != nil {
			log.Printf("[FoodTagsPageHandler] ❌ template food_tags_page: %v", err)
		}
	}
}

// FoodTagDeleteHandler удаляет метку + обновляет панель и облако (OOB)
func FoodTagDeleteHandler(tmpl *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err == nil && id > 0 {
			if err := DeleteTag(db, id); err != nil {
				log.Printf("[FoodTagDeleteHandler] ❌ DeleteTag: %v", err)
			}
		}

		tags, err := ListTags(db)
		if err != nil {
			http.Error(w, "Ошибка загрузки меток", http.StatusInternalServerError)
			return
		}

		_ = tmpl.ExecuteTemplate(w, "food_tags_page", map[string]any{
			"Tags": tags,
		})
	}
}
