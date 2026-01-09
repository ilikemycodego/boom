-- Все таблицы в одном файле

-- Траты
CREATE TABLE IF NOT EXISTS expenses (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  day         TEXT    NOT NULL,
  amount      INTEGER NOT NULL
);

-- Откладывания
CREATE TABLE IF NOT EXISTS deposits (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  day         TEXT    NOT NULL,
  amount      INTEGER NOT NULL
);

-- Цель (одна запись)
CREATE TABLE IF NOT EXISTS goals (
  id     INTEGER PRIMARY KEY CHECK (id = 1), -- Только одна запись
  amount INTEGER NOT NULL DEFAULT 0
);

-- Метки еды
CREATE TABLE IF NOT EXISTS food_tags (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  name       TEXT NOT NULL UNIQUE,
  created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Записи по дням (одна запись на дату)
CREATE TABLE IF NOT EXISTS food_entries (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  entry_date TEXT NOT NULL UNIQUE, -- YYYY-MM-DD
  created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Связь "запись дня" <-> "метка" (many-to-many)
CREATE TABLE IF NOT EXISTS food_entry_tags (
  entry_id INTEGER NOT NULL,
  tag_id   INTEGER NOT NULL,
  PRIMARY KEY (entry_id, tag_id),
  FOREIGN KEY (entry_id) REFERENCES food_entries(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id)   REFERENCES food_tags(id)   ON DELETE CASCADE
);

-- Индексы (не обязательно, но полезно)
CREATE INDEX IF NOT EXISTS idx_food_entry_tags_entry_id ON food_entry_tags(entry_id);
CREATE INDEX IF NOT EXISTS idx_food_entry_tags_tag_id   ON food_entry_tags(tag_id);