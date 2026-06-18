package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func Init() error {
	var err error

Database, err = sql.Open("sqlite3", "./backend/database/forum.db?_foreign_keys=on")
	if err != nil {
		return fmt.Errorf("can't open/create forum.db: %v", err)
	}

	if err := Database.Ping(); err != nil {
		return fmt.Errorf("can't connect to database: %v", err)
	}

schema, err := os.ReadFile("./backend/database/schema.sql")
	if err != nil {
		return fmt.Errorf("can't read schema: %v", err)
	}

	_, err = Database.Exec(string(schema)) // seeder
	if err != nil {
		return fmt.Errorf("schema execution failed: %v", err)
	}

	return nil
}
