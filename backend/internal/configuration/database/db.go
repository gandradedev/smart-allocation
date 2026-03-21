package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func New(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS assets (
			ticker          TEXT PRIMARY KEY,
			quantity        REAL NOT NULL DEFAULT 0,
			price           REAL NOT NULL DEFAULT 0,
			ceiling_price   REAL NOT NULL DEFAULT 0,
			target_percent  REAL NOT NULL DEFAULT 0,
			asset_type	 	TEXT NOT NULL,
			created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return err
	}

	return nil
}
