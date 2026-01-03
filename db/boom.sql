-- версия схемы
CREATE TABLE IF NOT EXISTS meta (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL
);

INSERT OR IGNORE INTO meta(key, value) VALUES ('schema_version', '1');

-- расходы
CREATE TABLE IF NOT EXISTS expenses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  occurred_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
  amount INTEGER NOT NULL CHECK(amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_expenses_occurred_at ON expenses(occurred_at);