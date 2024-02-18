package main

import (
	"dalennod/internal/db"
	"dalennod/internal/server"
	"dalennod/internal/setup"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var database *sql.DB = setup.CreateDB(setup.GetOS())
	server.Start(database)
	db.UserInput(database)
}
