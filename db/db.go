package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // ✅ регистрирует драйвер "sqlite"
)

// OpenDB открывает SQLite и создаёт таблицы
func OpenDB(path string) *sql.DB {
	dbConn, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatal("❌ open db:", err)
	}

	// Включаем нормальный режим блокировок (обычно полезно)
	if _, err := dbConn.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		log.Fatal("❌ pragma WAL:", err)
	}

	// Схема
	schema := `
CREATE TABLE IF NOT EXISTS expenses (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  day         TEXT    NOT NULL,
  amount      INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_expenses_day ON expenses(day);



-- ✅ отложил себе (deposit)
CREATE TABLE IF NOT EXISTS deposits (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  day         TEXT    NOT NULL,
  amount      INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_deposits_day ON deposits(day);

`
	if _, err := dbConn.Exec(schema); err != nil {
		log.Fatal("❌ schema:", err)
	}

	return dbConn
}
