package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
)

// Встраиваем SQL-файлы миграций в бинарник.
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

func Migrate(db *sql.DB) error {
	log.Println("🛠️ миграции: старт")

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Проверяем наличие ВСЕХ таблиц из схемы 001.sql
	allTablesExist, err := allTablesExistTx(tx, []string{
		"expenses",
		"deposits",
		"goals",
		"food_tags",
		"food_entries",
		"food_entry_tags",
	})
	if err != nil {
		return fmt.Errorf("check tables: %w", err)
	}

	// 1) Схема таблиц: исполняем только если НЕ все таблицы существуют
	if allTablesExist {
		log.Println("↩️ все таблицы уже есть — schema пропускаю (migrations/001.sql)")
	} else {
		if err := execSQLFile(tx, "migrations/001.sql"); err != nil {
			return err
		}
	}

	// 2) Seed: можно выполнять всегда (INSERT OR IGNORE безопасен)
	if err := execSQLFile(tx, "migrations/002.sql"); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	log.Println("✅ миграции: готово")
	return nil
}

func execSQLFile(tx *sql.Tx, path string) error {
	log.Printf("▶️ выполняю %s\n", path)

	b, err := migrationsFS.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	// SQLite нормально переваривает много стейтментов в одном Exec, разделенных ';'
	if _, err := tx.Exec(string(b)); err != nil {
		return fmt.Errorf("exec %s: %w", path, err)
	}

	log.Printf("✅ выполнено: %s\n", path)
	return nil
}

func allTablesExistTx(tx *sql.Tx, tables []string) (bool, error) {
	for _, t := range tables {
		exists, err := tableExistsTx(tx, t)
		if err != nil {
			return false, err
		}
		if !exists {
			return false, nil
		}
	}
	return true, nil

}

func tableExistsTx(tx *sql.Tx, tableName string) (bool, error) {
	var cnt int
	err := tx.QueryRow(
		`SELECT COUNT(1) FROM sqlite_master WHERE type='table' AND name=?;`,
		tableName,
	).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}
