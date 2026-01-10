package food

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Tag — метка еды
type Tag struct {
	ID       int
	Name     string
	Selected bool
}

// FoodPageData — данные для шаблона /food
type FoodPageData struct {
	Today        string
	Tags         []Tag
	SelectedTags []Tag
}

// todayStr возвращает дату в формате YYYY-MM-DD
func todayStr() string {
	return time.Now().Format("2006-01-02")
}

// ListTags возвращает все метки (по имени)
func ListTags(db *sql.DB) ([]Tag, error) {
	rows, err := db.Query(`SELECT id, name FROM food_tags ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// GetSelectedTagIDs возвращает выбранные метки для даты
func GetSelectedTagIDs(db *sql.DB, day string) (map[int]bool, error) {
	// Если записи нет — вернем пустой набор
	var entryID int
	err := db.QueryRow(`SELECT id FROM food_entries WHERE entry_date = ?`, day).Scan(&entryID)
	if err == sql.ErrNoRows {
		return map[int]bool{}, nil
	}
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`SELECT tag_id FROM food_entry_tags WHERE entry_id = ?`, entryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	selected := map[int]bool{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		selected[id] = true
	}
	return selected, rows.Err()
}

// SaveEntry сохраняет выбранные tagIDs как НОВУЮ запись (каждое сохранение = новая строка)
func SaveEntry(db *sql.DB, day string, tagIDs []int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// ✅ Всегда создаём НОВУЮ запись дня
	res, err := tx.Exec(`INSERT INTO food_entries(entry_date) VALUES (?)`, day)
	if err != nil {
		return err
	}

	entryID64, err := res.LastInsertId()
	if err != nil {
		return err
	}
	entryID := int(entryID64)

	// ✅ Привязываем метки к этой записи
	for _, tagID := range tagIDs {
		if tagID <= 0 {
			continue
		}
		// INSERT OR IGNORE — чтобы не падать, если вдруг один тег пришёл дважды
		if _, err := tx.Exec(
			`INSERT OR IGNORE INTO food_entry_tags(entry_id, tag_id) VALUES (?, ?)`,
			entryID, tagID,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// AddTag добавляет метку
func AddTag(db *sql.DB, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("пустое имя метки")
	}
	if len([]rune(name)) > 64 {
		return fmt.Errorf("слишком длинное имя метки (max 64)")
	}

	_, err := db.Exec(`INSERT INTO food_tags(name) VALUES (?)`, name)
	return err
}

// DeleteTag удаляет метку (связи удалятся каскадом)
func DeleteTag(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM food_tags WHERE id = ?`, id)
	return err
}
