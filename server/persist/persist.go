package persist

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
)

// ConnectAndMigrate opens DB connection and performs data migration on that DB
// used by runtime and integration tests
func ConnectAndMigrate(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return db, err
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return db, err
	}

	migration, err := migrate.NewWithDatabaseInstance("file://../migrations", "mysql", driver)
	if err != nil {
		return db, err
	}

	err = migration.Up()
	if err != nil {
		// should not fail if Database Schema is up to date or it is currently updated (and locked) by other instance
		if err == migrate.ErrNoChange || err == migrate.ErrLocked {
			fmt.Print("DB Migration: " + err.Error())
			return db, nil
		}
		return db, err
	}

	return db, nil
}
