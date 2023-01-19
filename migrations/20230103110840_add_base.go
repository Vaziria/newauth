package migrations

import (
	"database/sql"
	"log"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddBase, downAddBase)
}

func upAddBase(tx *sql.Tx) error {
	// This code is executed when the migration is applied.

	_, err := tx.Exec(`
		CREATE TABLE  users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(500) not NULL,
			username VARCHAR(500) NOT NULL UNIQUE,
			email TEXT not NULL UNIQUE,
			phone VARCHAR(100),
			password TEXT not NULL,
			suspended BOOLEAN DEFAULT FALSE,
			verified BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT now(),
			updated_at TIMESTAMP DEFAULT now(),
			last_reset TIMESTAMP DEFAULT now(),
			deleted_at TIMESTAMP NULL
		);
		
		CREATE TABLE teams (
			id INTEGER PRIMARY KEY,
			name VARCHAR(500) not NULL,
			description TEXT
		);

		CREATE TABLE permissions (
			id SERIAL PRIMARY KEY,
			key VARCHAR(1000) not NULL UNIQUE
		);

	`)

	if err != nil {
		return err
	}

	return nil
}

func downAddBase(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`
		DROP TABLE users;
		DROP TABLE teams;
		DROP TABLE permissions;
	`)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
