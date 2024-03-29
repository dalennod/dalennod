package setup

import (
	"database/sql"
	"fmt"
	"log"
)

type Bookmark struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Note        string `json:"note"`
	Keywords    string `json:"keywords"`
	BGroup      string `json:"bGroup"`
	Archived    bool   `json:"archive"`
	SnapshotURL string `json:"snapshotURL"`
	ThumbURL    string `json:"thumbURL"`
	Modified    string `json:"modified"`
}

const DB_FILENAME string = "default_user.db"

func CreateDB(dbSavePath string) *sql.DB {
	db, err := sql.Open("sqlite3", fmt.Sprint(dbSavePath+DB_FILENAME))
	if err != nil {
		log.Fatalln(err)
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
			thumbURL	TEXT,
			modified	DATETIME 	DEFAULT 		CURRENT_TIMESTAMP	NOT NULL
		);
	`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
