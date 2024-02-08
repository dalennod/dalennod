package setup

import (
	"dalennod/internal/logger"
	"database/sql"
	"fmt"
	"time"
)

const (
	DB_FILENAME = "default_user.db"
)

type Bookmark struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Note        string    `json:"note"`
	Keywords    string    `json:"keywords"`
	BGroup      string    `json:"bGroup"`
	Archived    bool      `json:"archive"`
	SnapshotURL string    `json:"snapshotURL"`
	Modified    time.Time `json:"modified"`
}

func CreateDB(dbSavePath string) *sql.DB {
	db, err := sql.Open("sqlite3", fmt.Sprint(dbSavePath+DB_FILENAME))
	if err != nil {
		logger.Error.Fatalln(err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS bookmarks (
			id 			INTEGER 	PRIMARY KEY 	NOT NULL,
			url 		TEXT 						NOT NULL,
			title 		TEXT,
			note 		TEXT,
			keywords	TEXT,
			bGroup		TEXT,
			archived	BOOLEAN		NOT NULL,
			snapshotURL	TEXT,
			modified	DATETIME 	DEFAULT 		CURRENT_TIMESTAMP	NOT NULL
		);
	`)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	execResult, err := stmt.Exec()
	if err != nil {
		logger.Error.Fatalln(err)
	}
	logger.Info.Println(execResult.RowsAffected())

	return db
}
