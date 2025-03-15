package setup

import (
    "database/sql"
    "log"
    "path/filepath"
)

type Bookmark struct {
    ID           int    `json:"id"`
    URL          string `json:"url"`
    Title        string `json:"title"`
    Note         string `json:"note"`
    Keywords     string `json:"keywords"`
    BmGroup      string `json:"bmGroup"`
    Archived     bool   `json:"archive"`
    SnapshotURL  string `json:"snapshotURL"`
    ThumbURL     string `json:"thumbURL"`
    ByteThumbURL []byte `json:"byteThumbURL"`
    Modified     string `json:"modified"`
}

const DB_FILENAME string = "dalennod.db"

func CreateDB(dbSavePath string) *sql.DB {
    db, err := sql.Open("sqlite3", filepath.Join(dbSavePath, DB_FILENAME)) // For CGo driver
    // db, err := sql.Open("sqlite", filepath.Join(dbSavePath, DB_FILENAME)) // For CGo-free driver
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
            bmGroup      TEXT,
            archived     BOOLEAN NOT NULL,
            snapshotURL  TEXT,
            thumbURL     TEXT,
            byteThumbURL BLOB,
            modified     DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
        );`)
    if err != nil {
        log.Fatalln("error when preparing database. ERROR:", err)
    }

    _, err = stmt.Exec()
    if err != nil {
        log.Fatalln("error creating database. ERROR:", err)
    }

    return db
}
