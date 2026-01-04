package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func OpenDB(dbPath string) *sql.DB {
	// 1. Создаем папку, если её нет (например, "./data")
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal("❌ не удалось создать папку для БД:", err)
	}

	// 2. Открываем соединение
	dbConn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("❌ open db:", err)
	}

	if _, err := dbConn.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		log.Fatal("❌ pragma WAL:", err)
	}

	// Схема (без изменений)
	schema := `
	CREATE TABLE IF NOT EXISTS expenses (
	  id          INTEGER PRIMARY KEY AUTOINCREMENT,
	  day         TEXT    NOT NULL,
	  amount      INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS deposits (
	  id          INTEGER PRIMARY KEY AUTOINCREMENT,
	  day         TEXT    NOT NULL,
	  amount      INTEGER NOT NULL
	);
	

CREATE TABLE IF NOT EXISTS goals (
  id    INTEGER PRIMARY KEY CHECK (id = 1), -- Только одна запись
  amount INTEGER NOT NULL DEFAULT 0
);
-- Инициализируем цель нулем, если записи еще нет
INSERT OR IGNORE INTO goals (id, amount) VALUES (1, 0);
	
	
	
	`

	if _, err := dbConn.Exec(schema); err != nil {
		log.Fatal("❌ schema:", err)
	}

	return dbConn
}
