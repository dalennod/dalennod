package db

import (
	"dalennod/internal/logger"
	"dalennod/internal/setup"
	"dalennod/internal/thumb_url"
	"database/sql"
	"fmt"
	"time"
)

func Add(database *sql.DB, bmStruct setup.Bookmark) {
	if bmStruct.ThumbURL == "" {
		var err error
		bmStruct.ThumbURL, err = thumb_url.GetPageThumb(bmStruct.URL)
		if err != nil {
			logger.Error.Println(err)
		}
	}

	stmt, err := database.Prepare("INSERT INTO bookmarks (url, title, note, keywords, bGroup, archived, snapshotURL, thumbURL) VALUES (?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	_, err = stmt.Exec(bmStruct.URL, bmStruct.Title, bmStruct.Note, bmStruct.Keywords, bmStruct.BGroup, bmStruct.Archived, bmStruct.SnapshotURL, bmStruct.ThumbURL)
	if err != nil {
		logger.Error.Fatalln(err)
	}
}

func Update(database *sql.DB, bmStruct setup.Bookmark, serverCall bool) {
	if !serverCall {
		bmStruct = updateCheck(database, bmStruct)
	}

	stmt, err := database.Prepare("UPDATE bookmarks SET url=(?), title=(?), note=(?), keywords=(?), bGroup=(?), archived=(?), snapshotURL=(?) WHERE id=(?);")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	_, err = stmt.Exec(bmStruct.URL, bmStruct.Title, bmStruct.Note, bmStruct.Keywords, bmStruct.BGroup, bmStruct.Archived, bmStruct.SnapshotURL, bmStruct.ID)
	if err != nil {
		logger.Error.Fatalln(err)
	}
}

func updateCheck(database *sql.DB, bmStruct setup.Bookmark) setup.Bookmark {
	oldBmData, err := ViewSingleRow(database, bmStruct.ID, true)
	if err != nil {
		logger.Error.Println(err)
	}

	if bmStruct.URL == "" {
		bmStruct.URL = oldBmData.URL
	}
	if bmStruct.Title == "" {
		bmStruct.Title = oldBmData.Title
	}
	if bmStruct.Note == "" {
		bmStruct.Note = oldBmData.Note
	}
	if bmStruct.Keywords == "" {
		bmStruct.Keywords = oldBmData.Keywords
	}
	if bmStruct.BGroup == "" {
		bmStruct.BGroup = oldBmData.BGroup
	}

	return bmStruct
}

func Remove(database *sql.DB, id int) {
	stmt, err := database.Prepare("DELETE FROM bookmarks WHERE id=(?);")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		logger.Error.Fatalln(err)
	}
}

// SELECT * FROM bookmarks ORDER by id DESC LIMIT 5 OFFSET 0;
// SELECT * FROM bookmarks ORDER by id DESC LIMIT 5 OFFSET 5;
// ... OFFSET 10;
// in 20 intervals for multiple pages
func ViewAll(database *sql.DB, serverCall bool) []setup.Bookmark {
	var results []setup.Bookmark
	var result setup.Bookmark
	var modified time.Time

	rows, err := database.Query("SELECT * FROM bookmarks ORDER BY id DESC;")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	for rows.Next() {
		result = setup.Bookmark{}
		rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &modified)
		if serverCall {
			results = append(results, setup.Bookmark{ID: result.ID, URL: result.URL, Title: result.Title, Note: result.Note, Keywords: result.Keywords, BGroup: result.BGroup, Archived: result.Archived, SnapshotURL: result.SnapshotURL, ThumbURL: result.ThumbURL, Modified: modified.Local().Format("2006-01-02 15:04:05")})
		} else {
			fmt.Printf("\t#%d\nTitle:\t\t%s\nURL:\t\t%s\nNote:\t\t%s\nKeywords:\t%s\nGroup:\t\t%s\nArchived?:\t%t\nArchive URL:\t%s\nModified:\t%v\n\n", result.ID, result.Title, result.URL, result.Note, result.Keywords, result.BGroup, result.Archived, result.SnapshotURL, modified.Local().Format("2006-01-02 15:04:05"))
		}
	}
	return results
}

func ViewAllWhere(database *sql.DB, keyword string) []setup.Bookmark {
	var results []setup.Bookmark
	var result setup.Bookmark
	var modified time.Time

	if keyword == "" {
		return results
	}
	keyword = "%" + keyword + "%"

	stmt, err := database.Prepare("SELECT * FROM bookmarks WHERE keywords LIKE (?) or bGroup LIKE (?) or note LIKE (?) or title LIKE (?) ORDER BY id DESC;")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	execRes, err := stmt.Query(keyword, keyword, keyword, keyword)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	for execRes.Next() {
		result = setup.Bookmark{}
		execRes.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &modified)
		results = append(results, setup.Bookmark{ID: result.ID, URL: result.URL, Title: result.Title, Note: result.Note, Keywords: result.Keywords, BGroup: result.BGroup, Archived: result.Archived, SnapshotURL: result.SnapshotURL, ThumbURL: result.ThumbURL, Modified: modified.Local().Format("2006-01-02 15:04:05")})
	}

	return results
}

func ViewSingleRow(database *sql.DB, id int, serverCall bool) (setup.Bookmark, error) {
	var rowResult setup.Bookmark
	var modified time.Time

	stmt, err := database.Prepare("SELECT * FROM bookmarks WHERE id=(?);")
	if err != nil {
		return rowResult, err
	}

	execRes, err := stmt.Query(id)
	if err != nil {
		return rowResult, err
	}

	for execRes.Next() {
		err = execRes.Scan(&rowResult.ID, &rowResult.URL, &rowResult.Title, &rowResult.Note, &rowResult.Keywords, &rowResult.BGroup, &rowResult.Archived, &rowResult.SnapshotURL, &rowResult.ThumbURL, &modified)
		if err != nil {
			return rowResult, err
		}
	}

	if rowResult.URL == "" {
		return rowResult, fmt.Errorf("ID does not exist")
	}

	if !serverCall {
		fmt.Printf("\t#%d\nTitle:\t\t%s\nURL:\t\t%s\nNote:\t\t%s\nKeywords:\t%s\nGroup:\t\t%s\nArchived?:\t%t\nArchive URL:\t%s\nModified:\t%v\n", rowResult.ID, rowResult.Title, rowResult.URL, rowResult.Note, rowResult.Keywords, rowResult.BGroup, rowResult.Archived, rowResult.SnapshotURL, modified.Local().Format("2006-01-02 15:04:05"))
		return rowResult, nil
	}

	return rowResult, nil
}

func SearchByUrl(database *sql.DB, searchUrl string) (setup.Bookmark, error) {
	var foundBookmark setup.Bookmark

	stmt, err := database.Prepare("SELECT * FROM bookmarks WHERE url=(?);")
	if err != nil {
		return foundBookmark, err
	}
	execRes, err := stmt.Query(searchUrl)
	if err != nil {
		return foundBookmark, err
	}

	for execRes.Next() {
		err = execRes.Scan(&foundBookmark.ID, &foundBookmark.URL, &foundBookmark.Title, &foundBookmark.Note, &foundBookmark.Keywords, &foundBookmark.BGroup, &foundBookmark.Archived, &foundBookmark.SnapshotURL, &foundBookmark.ThumbURL, &foundBookmark.Modified)
		if err != nil {
			return foundBookmark, err
		}
	}

	return foundBookmark, nil
}
