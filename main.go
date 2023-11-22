package main

import (
	"dalennod/internal/db"
	"dalennod/internal/server"
	"dalennod/internal/setup"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// var dbSavePath string = setup.GetOS()
	// db := setup.CreateDB(dbSavePath)

	var database *sql.DB = setup.CreateDB(setup.GetOS())
	go server.Start(database)
	db.UserInput(database)

}
