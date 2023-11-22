package setup

import (
	"dalennod/internal/logger"
	"database/sql"
	"fmt"
)

const (
	DB_FILENAME = "default_user.db"
)

func CreateDB(dbSavePath string) *sql.DB {
	db, err := sql.Open("sqlite3", fmt.Sprint(dbSavePath+DB_FILENAME))
	if err != nil {
		logger.Error.Fatalln(err)
	}

	// add column for direct archived snapshot
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS bookmarks (
			id 			INTEGER 	PRIMARY KEY 	NOT NULL,
			url 		TEXT 						NOT NULL,
			title 		TEXT,
			note 		TEXT,
			keywords	TEXT,
			bGroup		TEXT,
			archived	BOOLEAN		NOT NULL,
			modified	DATETIME 	DEFAULT 		CURRENT_TIMESTAMP	NOT NULL
		);
	`)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	// log.Println(stmt)

	execResult, err := stmt.Exec()
	if err != nil {
		logger.Error.Fatalln(err)
	}
	logger.Info.Println(execResult.RowsAffected())

	return db
}
