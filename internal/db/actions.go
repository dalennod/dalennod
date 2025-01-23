package db

import (
	"dalennod/internal/logger"
	"dalennod/internal/setup"
	"dalennod/internal/thumb_url"
	"database/sql"
	"fmt"
	"time"
)

const PAGE_UPDATE_LIMIT = 60

func Add(database *sql.DB, bmStruct setup.Bookmark) {
	if bmStruct.ThumbURL == "" || len(bmStruct.ByteThumbURL) == 0 {
		var err error
		bmStruct.ThumbURL, bmStruct.ByteThumbURL, err = thumb_url.GetPageThumb(bmStruct.URL)
		if err != nil {
			logger.Error.Println(err)
		}
	}

	stmt, err := database.Prepare("INSERT INTO bookmarks (url, title, note, keywords, bmGroup, archived, snapshotURL, thumbURL, byteThumbURL) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	_, err = stmt.Exec(bmStruct.URL, bmStruct.Title, bmStruct.Note, bmStruct.Keywords, bmStruct.BmGroup, bmStruct.Archived, bmStruct.SnapshotURL, bmStruct.ThumbURL, bmStruct.ByteThumbURL)
	if err != nil {
		logger.Error.Fatalln(err)
	}
}

func Update(database *sql.DB, bmStruct setup.Bookmark, serverCall bool) {
	if !serverCall {
		bmStruct = updateCheck(database, bmStruct)
	}

	stmt, err := database.Prepare("UPDATE bookmarks SET url=(?), title=(?), note=(?), keywords=(?), bmGroup=(?), archived=(?), snapshotURL=(?), thumbURL=(?), byteThumbURL=(?) WHERE id=(?);")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	_, err = stmt.Exec(bmStruct.URL, bmStruct.Title, bmStruct.Note, bmStruct.Keywords, bmStruct.BmGroup, bmStruct.Archived, bmStruct.SnapshotURL, bmStruct.ThumbURL, bmStruct.ByteThumbURL, bmStruct.ID)
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
	if bmStruct.BmGroup == "" {
		bmStruct.BmGroup = oldBmData.BmGroup
	}

	return bmStruct
}

func RefetchThumbnail(database *sql.DB, id int, thumbnail []byte) error {
	bmStruct, err := ViewSingleRow(database, id, true)
	if err != nil {
		logger.Error.Println("error getting single row. ERROR:", err)
		return err
	}

	if thumbnail == nil {
		bmStruct.ThumbURL, bmStruct.ByteThumbURL, err = thumb_url.GetPageThumb(bmStruct.URL)
		if err != nil {
			logger.Error.Println("error getting thumbnail. ERROR:", err)
			return err
		}
	} else {
		bmStruct.ByteThumbURL = thumbnail
	}

	Update(database, bmStruct, true)
	return nil
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

func TotalPageCount(database *sql.DB) int {
	rows, err := database.Query("SELECT COUNT(*) FROM bookmarks;")
	if err != nil {
		logger.Error.Println("error getting total page count from database. ERROR:", err)
	}

	var pageCount int
	for rows.Next() {
		rows.Scan(&pageCount)
	}
	pageCount = pageCount / PAGE_UPDATE_LIMIT
	return pageCount
}

func ViewAllWebUI(database *sql.DB, pageNo int) []setup.Bookmark {
	var results []setup.Bookmark
	var result setup.Bookmark
	var modified time.Time

	stmt, err := database.Prepare("SELECT * FROM bookmarks ORDER BY id DESC LIMIT (?) OFFSET (?);")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	pageOffset := pageNo * PAGE_UPDATE_LIMIT

	rows, err := stmt.Query(PAGE_UPDATE_LIMIT, pageOffset)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	for rows.Next() {
		result = setup.Bookmark{}
		rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
		results = append(results, setup.Bookmark{ID: result.ID, URL: result.URL, Title: result.Title, Note: result.Note, Keywords: result.Keywords, BmGroup: result.BmGroup, Archived: result.Archived, SnapshotURL: result.SnapshotURL, ThumbURL: result.ThumbURL, ByteThumbURL: result.ByteThumbURL, Modified: modified.Local().Format("2006-01-02 15:04:05")})
	}
	return results
}

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
		rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
		if serverCall {
			results = append(results, setup.Bookmark{ID: result.ID, URL: result.URL, Title: result.Title, Note: result.Note, Keywords: result.Keywords, BmGroup: result.BmGroup, Archived: result.Archived, SnapshotURL: result.SnapshotURL, ThumbURL: result.ThumbURL, ByteThumbURL: result.ByteThumbURL, Modified: modified.Local().Format("2006-01-02 15:04:05")})
		} else {
			fmt.Printf("\t#%d\nTitle:\t\t%s\nURL:\t\t%s\nNote:\t\t%s\nKeywords:\t%s\nGroup:\t\t%s\nArchived?:\t%t\nArchive URL:\t%s\nModified:\t%v\n\n", result.ID, result.Title, result.URL, result.Note, result.Keywords, result.BmGroup, result.Archived, result.SnapshotURL, modified.Local().Format("2006-01-02 15:04:05"))
		}
	}
	return results
}

func BackupViewAll(database *sql.DB) []setup.Bookmark {
	var results []setup.Bookmark
	var result setup.Bookmark
	var modified time.Time

	rows, err := database.Query("SELECT * FROM bookmarks ORDER BY id;")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	for rows.Next() {
		result = setup.Bookmark{}
		rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
		results = append(results, setup.Bookmark{ID: result.ID, URL: result.URL, Title: result.Title, Note: result.Note, Keywords: result.Keywords, BmGroup: result.BmGroup, Archived: result.Archived, SnapshotURL: result.SnapshotURL, ThumbURL: result.ThumbURL, ByteThumbURL: result.ByteThumbURL, Modified: modified.Local().Format("2006-01-02 15:04:05")})
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

	stmt, err := database.Prepare("SELECT * FROM bookmarks WHERE keywords LIKE (?) or bmGroup LIKE (?) or note LIKE (?) or title LIKE (?) or url LIKE (?) ORDER BY id DESC;")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	execRes, err := stmt.Query(keyword, keyword, keyword, keyword, keyword)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	for execRes.Next() {
		result = setup.Bookmark{}
		execRes.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
		results = append(results, setup.Bookmark{ID: result.ID, URL: result.URL, Title: result.Title, Note: result.Note, Keywords: result.Keywords, BmGroup: result.BmGroup, Archived: result.Archived, SnapshotURL: result.SnapshotURL, ThumbURL: result.ThumbURL, ByteThumbURL: result.ByteThumbURL, Modified: modified.Local().Format("2006-01-02 15:04:05")})
	}

	return results
}

func ViewAllWhereKeyword(database *sql.DB, keyword string) []setup.Bookmark {
	var results []setup.Bookmark
	var result setup.Bookmark
	var modified time.Time

	keyword = "%" + keyword + "%"

	stmt, err := database.Prepare("SELECT * FROM bookmarks WHERE keywords LIKE (?) ORDER BY id DESC;")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	execRes, err := stmt.Query(keyword)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	for execRes.Next() {
		result = setup.Bookmark{}
		execRes.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
		results = append(results, setup.Bookmark{ID: result.ID, URL: result.URL, Title: result.Title, Note: result.Note, Keywords: result.Keywords, BmGroup: result.BmGroup, Archived: result.Archived, SnapshotURL: result.SnapshotURL, ThumbURL: result.ThumbURL, ByteThumbURL: result.ByteThumbURL, Modified: modified.Local().Format("2006-01-02 15:04:05")})
	}

	return results
}

func ViewAllWhereGroup(database *sql.DB, keyword string) []setup.Bookmark {
	var results []setup.Bookmark
	var result setup.Bookmark
	var modified time.Time

	keyword = "%" + keyword + "%"

	stmt, err := database.Prepare("SELECT * FROM bookmarks WHERE bmGroup LIKE (?) ORDER BY id DESC;")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	execRes, err := stmt.Query(keyword)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	for execRes.Next() {
		result = setup.Bookmark{}
		execRes.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
		results = append(results, setup.Bookmark{ID: result.ID, URL: result.URL, Title: result.Title, Note: result.Note, Keywords: result.Keywords, BmGroup: result.BmGroup, Archived: result.Archived, SnapshotURL: result.SnapshotURL, ThumbURL: result.ThumbURL, ByteThumbURL: result.ByteThumbURL, Modified: modified.Local().Format("2006-01-02 15:04:05")})
	}

	return results
}

func ViewAllWhereHostname(database *sql.DB, hostname string) []setup.Bookmark {
	var results []setup.Bookmark
	var result setup.Bookmark
	var modified time.Time

	hostname = "%" + hostname + "%"

	stmt, err := database.Prepare("SELECT * FROM bookmarks WHERE url LIKE (?) ORDER BY id DESC;")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	execRes, err := stmt.Query(hostname)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	for execRes.Next() {
		result = setup.Bookmark{}
		execRes.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
		results = append(results, setup.Bookmark{ID: result.ID, URL: result.URL, Title: result.Title, Note: result.Note, Keywords: result.Keywords, BmGroup: result.BmGroup, Archived: result.Archived, SnapshotURL: result.SnapshotURL, ThumbURL: result.ThumbURL, ByteThumbURL: result.ByteThumbURL, Modified: modified.Local().Format("2006-01-02 15:04:05")})
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
		err = execRes.Scan(&rowResult.ID, &rowResult.URL, &rowResult.Title, &rowResult.Note, &rowResult.Keywords, &rowResult.BmGroup, &rowResult.Archived, &rowResult.SnapshotURL, &rowResult.ThumbURL, &rowResult.ByteThumbURL, &modified)
		if err != nil {
			return rowResult, err
		}
	}

	if rowResult.URL == "" {
		return rowResult, fmt.Errorf("ID does not exist")
	}

	if !serverCall {
		fmt.Printf("\t#%d\nTitle:\t\t%s\nURL:\t\t%s\nNote:\t\t%s\nKeywords:\t%s\nGroup:\t\t%s\nArchived?:\t%t\nArchive URL:\t%s\nModified:\t%v\n", rowResult.ID, rowResult.Title, rowResult.URL, rowResult.Note, rowResult.Keywords, rowResult.BmGroup, rowResult.Archived, rowResult.SnapshotURL, modified.Local().Format("2006-01-02 15:04:05"))
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
		var err error = execRes.Scan(&foundBookmark.ID, &foundBookmark.URL, &foundBookmark.Title, &foundBookmark.Note, &foundBookmark.Keywords, &foundBookmark.BmGroup, &foundBookmark.Archived, &foundBookmark.SnapshotURL, &foundBookmark.ThumbURL, &foundBookmark.ByteThumbURL, &foundBookmark.Modified)
		if err != nil {
			return foundBookmark, err
		}
	}

	return foundBookmark, nil
}

func GetAllGroups(database *sql.DB) ([]string, error) {
	var allGroups []string

	rows, err := database.Query("SELECT DISTINCT bmGroup FROM bookmarks ORDER BY id DESC;")
	if err != nil {
		return allGroups, err
	}

	var bmGroup string
	for rows.Next() {
		var err error = rows.Scan(&bmGroup)
		if err != nil {
			logger.Error.Printf("error when scanning row. ERROR: %v", err)
		}
		if bmGroup != "" {
			allGroups = append(allGroups, bmGroup)
		}
	}

	return allGroups, nil
}
