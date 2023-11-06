package model

import (
	"database/sql"
	"user-track/global"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "wordcards.db")
	if err != nil {
		return err
	}
	return nil
}

// Migrate Model
func AutoMigrateAll() {
	migrateErr := global.GormDb.AutoMigrate(
	// &Stock{},
	)
	if migrateErr != nil {
		panic("database migrate failed")
	}
}
