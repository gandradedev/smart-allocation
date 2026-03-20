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
			created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return err
	}

	// Drop aporte_ajustado if it exists (now a calculated field).
	// Ignore errors — column may not exist on fresh installs.
	db.Exec(`ALTER TABLE assets DROP COLUMN aporte_ajustado`)

	// Rename Portuguese columns to English on existing databases.
	// Ignore errors — columns may already be renamed or not exist.
	db.Exec(`ALTER TABLE assets RENAME COLUMN qtd_acoes TO quantity`)
	db.Exec(`ALTER TABLE assets RENAME COLUMN preco TO price`)
	db.Exec(`ALTER TABLE assets RENAME COLUMN preco_teto TO ceiling_price`)
	db.Exec(`ALTER TABLE assets RENAME COLUMN objetivo_percent TO target_percent`)

	return nil
}
