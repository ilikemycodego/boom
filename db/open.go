package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func OpenDB(dbPath string) *sql.DB {
	// 1) Создаем папку, если её нет (например, "./data")
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal("❌ не удалось создать папку для БД:", err)
	}

	// 2) Открываем соединение
	dbConn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("❌ open db:", err)
	}

	// 3) PRAGMA-настройки (важно включить foreign_keys для ON DELETE CASCADE)
	if _, err := dbConn.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		log.Fatal("❌ pragma WAL:", err)
	}
	if _, err := dbConn.Exec(`PRAGMA foreign_keys=ON;`); err != nil {
		log.Fatal("❌ pragma foreign_keys:", err)
	}

	// Для SQLite обычно полезно держать 1 соединение
	dbConn.SetMaxOpenConns(1)
	dbConn.SetMaxIdleConns(1)

	// Чтобы не падать "database is locked" при конкуренции
	if _, err := dbConn.Exec(`PRAGMA busy_timeout=5000;`); err != nil {
		log.Fatal("❌ pragma busy_timeout:", err)
	}

	// Проверка соединения сразу
	if err := dbConn.Ping(); err != nil {
		log.Fatal("❌ ping db:", err)
	}

	// 4) Миграции из SQL-файлов
	if err := Migrate(dbConn); err != nil {
		log.Fatal("❌ migrate:", err)
	}

	return dbConn
}
