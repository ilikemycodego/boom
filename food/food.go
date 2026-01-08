package food

import (
	"html/template"
	"log"
	"net/http"
)

func FoodHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// --- рендер шаблона ---
		if err := tmpl.ExecuteTemplate(w, "food", nil); err != nil {
			log.Printf("[FoodHandler] ❌ Ошибка шаблона: %v", err)
			http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
			return
		}

		log.Println("[FoodHandler] ✅ FoodHandler отрендерена")
	}
}
