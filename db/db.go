package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() *sql.DB {
	var err error
	DB, err = sql.Open("sqlite", "api.db")

	if err != nil {
		panic("Cloud not connect to database")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()

	return DB
}

func createTables() {
	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)
	`

	_, err := DB.Exec(createUsersTable)

	if err != nil {
		fmt.Println(err)
		panic("Cloud Not Create users Table")
	}

	createEventsTable := `
		CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			location TEXT NOT NULL,
			date_time DATETIME NOT NULL,
			user_id TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)
	`

	_, err = DB.Exec(createEventsTable)
	if err != nil {
		fmt.Println(err)
		panic("Cloud NOT Create Events Table")
	}

}
