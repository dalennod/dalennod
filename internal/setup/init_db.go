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
	BmGroup     string `json:"bmGroup"`
	Archived    bool   `json:"archive"`
	SnapshotURL string `json:"snapshotURL"`
	ThumbURL    string `json:"thumbURL"`
	B64ThumbURL string `json:"b64ThumbURL"`
	Modified    string `json:"modified"`
}

const DB_FILENAME string = "dalennod.db"

func CreateDB(dbSavePath string) *sql.DB {
	db, err := sql.Open("sqlite3", fmt.Sprint(dbSavePath+DB_FILENAME))
	if err != nil {
		log.Fatalln(err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS bookmarks (
			id INTEGER PRIMARY KEY,
			url TEXT NOT NULL,
			title TEXT,
			note TEXT,
			keywords TEXT,
			bmGroup TEXT,
			archived BOOLEAN NOT NULL,
			snapshotURL TEXT,
			thumbURL TEXT,
			b64ThumbURL TEXT,
			modified DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
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
