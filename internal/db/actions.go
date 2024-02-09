package db

import (
	"dalennod/internal/logger"
	"dalennod/internal/setup"
	"database/sql"
	"fmt"
	"time"
)

var (
	id          int
	url         string
	title       string
	note        string
	keywords    string
	bGroup      string
	archived    bool
	snapshotURL string
	modified    time.Time
)

func Add(database *sql.DB, url, title, note, keywords, bGroup string, archive bool, snapshotURL string) {
	stmt, err := database.Prepare("INSERT INTO bookmarks (url, title, note, keywords, bGroup, archived, snapshotURL) VALUES (?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	execResult, err := stmt.Exec(url, title, note, keywords, bGroup, archive, snapshotURL)
	if err != nil {
		logger.Error.Fatalln(err)
	}
	logger.Info.Println(execResult.RowsAffected())
}

func Update(database *sql.DB, url, title, note, keywords, bGroup string, id int, archive, s bool, snapshotURL string) {
	if !s {
		url, title, note, keywords, bGroup = updateCheck(database, url, title, note, keywords, bGroup, id)
	}

	logger.Info.Printf("url %s, title %s, note %s, keywords %s, bGroup %s, id %d, archive %t, s %t, snapshotURL %s\n", url, title, note, keywords, bGroup, id, archive, s, snapshotURL)

	stmt, err := database.Prepare("UPDATE bookmarks SET url=?, title=?, note=?, keywords=?, bGroup=?, archived=?, snapshotURL=? WHERE id=?;")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	exec, err := stmt.Exec(url, title, note, keywords, bGroup, archive, snapshotURL, id)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	logger.Info.Println(exec.RowsAffected())
}

func updateCheck(database *sql.DB, url, title, note, keywords, bGroup string, id int) (string, string, string, string, string) {
	var (
		oldURL, oldTitle, oldNote, oldKeywords, oldBGroup string
	)
	oldURL, oldTitle, oldNote, oldKeywords, oldBGroup, _ = ViewSingleRow(database, id, true)

	if url == "" {
		url = oldURL
	}
	if title == "" {
		title = oldTitle
	}
	if note == "" {
		note = oldNote
	}
	if keywords == "" {
		keywords = oldKeywords
	}
	if bGroup == "" {
		bGroup = oldBGroup
	}

	return url, title, note, keywords, bGroup
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

func RemoveAll(database *sql.DB) {
	_, err := database.Query("DELETE FROM bookmarks;")
	if err != nil {
		logger.Error.Println(err)
	}
}

func ViewAll(database *sql.DB, s bool) []setup.Bookmark {
	var result []setup.Bookmark

	rows, err := database.Query("SELECT * FROM bookmarks;")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	for rows.Next() {
		rows.Scan(&id, &url, &title, &note, &keywords, &bGroup, &archived, &snapshotURL, &modified)
		if s {
			result = append(result, setup.Bookmark{ID: id, URL: url, Title: title, Note: note, Keywords: keywords, BGroup: bGroup, Archived: archived, SnapshotURL: snapshotURL, Modified: modified})
		} else {
			fmt.Printf("%d : %s, %s, %s, %s, %s, %t, %s, %v\n", id, url, title, note, keywords, bGroup, archived, snapshotURL, modified)
		}
	}
	return result
}

func ViewAllWhere(database *sql.DB, keyword string) []setup.Bookmark {
	var result []setup.Bookmark
	if keyword == "" {
		return result
	}
	keyword = "%" + keyword + "%"

	stmt, err := database.Prepare("SELECT * FROM bookmarks WHERE keywords LIKE (?) or bGroup LIKE (?) or note LIKE (?) or title LIKE (?);")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	execRes, err := stmt.Query(keyword, keyword, keyword, keyword)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	for execRes.Next() {
		execRes.Scan(&id, &url, &title, &note, &keywords, &bGroup, &archived, &snapshotURL, &modified)
		result = append(result, setup.Bookmark{ID: id, URL: url, Title: title, Note: note, Keywords: keywords, BGroup: bGroup, Archived: archived, SnapshotURL: snapshotURL, Modified: modified})
	}

	return result
}

func ViewSingleRow(database *sql.DB, id int, s bool) (string, string, string, string, string, bool) {
	// rows, err := database.Query(fmt.Sprintf("SELECT * FROM bookmarks WHERE id=%d;", id))
	// if err != nil {
	// 	logger.Error.Fatalln(err)
	// }

	stmt, err := database.Prepare("SELECT * FROM bookmarks WHERE id=(?);")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	execRes, err := stmt.Query(id)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	for execRes.Next() {
		execRes.Scan(&id, &url, &title, &note, &keywords, &bGroup, &archived, &snapshotURL, &modified)
		if !s {
			fmt.Printf("%d : %s, %s, %s, %s, %s, %t, %s, %v\n", id, url, title, note, keywords, bGroup, archived, snapshotURL, modified)
		}
	}
	return url, title, note, keywords, bGroup, archived
}
