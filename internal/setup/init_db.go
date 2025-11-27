package setup

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"dalennod/internal/constants"
)

type Bookmark struct {
	ID          int64  `json:"id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Note        string `json:"note"`
	Keywords    string `json:"keywords"`
	Category    string `json:"category"`
	Archived    bool   `json:"archive"`
	SnapshotURL string `json:"snapshotURL"`
	ThumbURL    string `json:"thumbURL"`
	Modified    string `json:"modified"`
}

type RecentInteractions struct {
	Bookmarks    []Bookmark `json:"bookmarks"`
	ID           int64      `json:"id"`
	BookmarkID   int64      `json:"bookmarkID"`
	LastAccessed string     `json:"lastAccessed"`
}

func CreateDB(dbSavePath string) *sql.DB {
	dbWithConnStrings := fmt.Sprintf("%s?_foreign_keys=true", filepath.Join(dbSavePath, constants.DB_FILENAME))
	db, err := sql.Open("sqlite3", dbWithConnStrings)
	if err != nil {
		log.Fatalln("error opening database. ERROR:", err)
	}

	stmt, err := db.Prepare(`
        CREATE TABLE IF NOT EXISTS bookmarks (
            id           INTEGER PRIMARY KEY,
            url          TEXT NOT NULL,
            title        TEXT,
            note         TEXT,
            keywords     TEXT,
            category     TEXT,
            archived     BOOLEAN NOT NULL,
            snapshotURL  TEXT,
            thumbURL     TEXT,
            modified     DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
        );`)
	if err != nil {
		log.Fatalln("error when preparing database. ERROR:", err)
	}

	if _, err = stmt.Exec(); err != nil {
		log.Fatalln("error creating 'bookmarks' table. ERROR:", err)
	}

	stmt, err = db.Prepare(`
        CREATE TABLE IF NOT EXISTS recents (
            id           INTEGER PRIMARY KEY,
            bookmarkID   INTEGER NOT NULL UNIQUE,
            lastAccessed DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
            FOREIGN KEY  (bookmarkID) REFERENCES bookmarks(id) ON DELETE CASCADE
        );`)
	if err != nil {
		log.Fatalln("error preparing database. ERROR:", err)
	}

	if _, err = stmt.Exec(); err != nil {
		log.Fatalln("error creating 'recents' table. ERROR:", err)
	}

	return db
}
