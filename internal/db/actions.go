package db

import (
	"dalennod/internal/logger"
	"database/sql"
	"fmt"
	"time"
)

var (
	id       int
	url      string
	title    string
	note     string
	keywords string
	bGroup   string
	archived int
	modified time.Time
)

type Bookmark struct {
	ID       int
	URL      string
	Title    string
	Note     string
	Keywords string
	BGroup   string
	Archived int
	Modified time.Time
}

func Add(database *sql.DB, url, title, note, keywords, bGroup string, archive int) {
	stmt, err := database.Prepare("INSERT INTO bookmarks (url, title, note, keywords, bGroup, archived) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	execResult, err := stmt.Exec(url, title, note, keywords, bGroup, archive)
	if err != nil {
		logger.Error.Fatalln(err)
	}
	logger.Info.Println(execResult.RowsAffected())
}

func Remove(database *sql.DB, id int) {
	stmt, err := database.Prepare("DELETE FROM bookmarks WHERE id=(?);")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	execResult, err := stmt.Exec(id)
	if err != nil {
		logger.Error.Fatalln(err)
	}
	logger.Info.Println(execResult.RowsAffected())
}

func ViewAll(database *sql.DB, o string) []Bookmark {
	var result []Bookmark

	rows, err := database.Query("SELECT * FROM bookmarks;")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	for rows.Next() {
		rows.Scan(&id, &url, &title, &note, &keywords, &bGroup, &archived, &modified)
		if o != "s" {
			fmt.Printf("%d : %s, %s, %s, %s, %s, %d, %v\n", id, url, title, note, keywords, bGroup, archived, modified)
		}
		result = append(result, Bookmark{id, url, title, note, keywords, bGroup, archived, modified})
	}
	// logger.Info.Println(result)
	return result
}

func ViewSingleRow(database *sql.DB, id int) {
	rows, err := database.Query(fmt.Sprintf("SELECT * FROM bookmarks WHERE id=%d;", id))
	if err != nil {
		logger.Error.Fatalln(err)
	}

	for rows.Next() {
		rows.Scan(&id, &url, &title, &note, &keywords, &bGroup, &archived, &modified)
		fmt.Printf("%d : %s, %s, %s, %s, %s, %d, %v\n", id, url, title, note, keywords, bGroup, archived, modified)
	}
}
