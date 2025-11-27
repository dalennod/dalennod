package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"dalennod/internal/constants"
	"dalennod/internal/setup"
	"dalennod/internal/thumb_url"
)

func Add(database *sql.DB, bkmStruct setup.Bookmark) {
	if bkmStruct.ThumbURL == "" {
		var err error
		bkmStruct.ThumbURL, err = thumb_url.GetPageThumb(bkmStruct.URL)
		if err != nil {
			log.Printf("WARN: could not get webpage thumbnail for: %s: %v:\n", bkmStruct.URL, err)
		}
	}

	stmt, err := database.Prepare("INSERT INTO bookmarks (url, title, note, keywords, category, archived, snapshotURL, thumbURL) VALUES (?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		log.Printf("WARN: preparing database statement to create bookmark: %s: %v\n", bkmStruct.URL, err)
		return
	}

	execResult, err := stmt.Exec(bkmStruct.URL, bkmStruct.Title, bkmStruct.Note, bkmStruct.Keywords, bkmStruct.Category, bkmStruct.Archived, bkmStruct.SnapshotURL, bkmStruct.ThumbURL)
	if err != nil {
		log.Printf("WARN: error executing database statement to create bookmark: %s: %v\n", bkmStruct.URL, err)
		return
	}

	lastInsertID, err := execResult.LastInsertId()
	if err != nil {
		log.Println("WARN: could not get last insert ID:", err)
		return
	}

	go saveThumbLocally(lastInsertID, bkmStruct.ThumbURL)
}

// serverCall boolean explanation:
// When updating from CLI, newline/empty input means to retain old data.
// But, when updating from Web UI or extension the empty input means no
// input. So, if updating from CLI the old data needs to be retrieved to
// replace the empty input, thus if not serverCall then get old data and
// replace the empty input with old data.
func Update(database *sql.DB, bkmStruct setup.Bookmark, serverCall bool) {
	if !serverCall {
		bkmStruct = updateCheck(database, bkmStruct)
	}

	if bkmStruct.ThumbURL == "" {
		var err error
		bkmStruct.ThumbURL, err = thumb_url.GetPageThumb(bkmStruct.URL)
		if err != nil {
			log.Printf("WARN: could not get webpage thumbnail for: %d: %v:\n", bkmStruct.ID, err)
		}
	}

	stmt, err := database.Prepare("UPDATE bookmarks SET url=(?), title=(?), note=(?), keywords=(?), category=(?), archived=(?), snapshotURL=(?), thumbURL=(?), modified=CURRENT_TIMESTAMP WHERE id=(?);")
	if err != nil {
		log.Printf("WARN: preparing database statement to update bookmark: %d: %v\n", bkmStruct.ID, err)
		return
	}

	_, err = stmt.Exec(bkmStruct.URL, bkmStruct.Title, bkmStruct.Note, bkmStruct.Keywords, bkmStruct.Category, bkmStruct.Archived, bkmStruct.SnapshotURL, bkmStruct.ThumbURL, bkmStruct.ID)
	if err != nil {
		log.Printf("WARN: preparing database statement to update bookmark: %d: %v\n", bkmStruct.ID, err)
		return
	}

	go saveThumbLocally(bkmStruct.ID, bkmStruct.ThumbURL)
}

func updateCheck(database *sql.DB, bkmStruct setup.Bookmark) setup.Bookmark {
	oldBKMData, err := ViewSingleRow(database, bkmStruct.ID)
	if err != nil {
		log.Printf("WARN: could not get bookmark row data at ID %d: %v\n", bkmStruct.ID, err)
		return bkmStruct
	}

	if bkmStruct.URL == "" {
		bkmStruct.URL = oldBKMData.URL
	}
	if bkmStruct.Title == "" {
		bkmStruct.Title = oldBKMData.Title
	}
	if bkmStruct.Note == "" {
		bkmStruct.Note = oldBKMData.Note
	}
	if bkmStruct.Keywords == "" {
		bkmStruct.Keywords = oldBKMData.Keywords
	}
	if bkmStruct.Category == "" {
		bkmStruct.Category = oldBKMData.Category
	}

	return bkmStruct
}

func RefetchThumbnail(database *sql.DB, id int64, thumbnail []byte) error {
	bkmStruct, err := ViewSingleRow(database, id)
	if err != nil {
		log.Printf("WARN: could not get bookmark row data at ID %d: %v\n", id, err)
		return err
	}

	if thumbnail == nil {
		bkmStruct.ThumbURL, err = thumb_url.GetPageThumb(bkmStruct.URL)
		if err != nil {
			log.Println("WARN: could not get webpage thumbnail:", err)
			return err
		}
	}

	Update(database, bkmStruct, true)
	return nil
}

func Remove(database *sql.DB, id int64) {
	stmt, err := database.Prepare("DELETE FROM bookmarks WHERE id=(?);")
	if err != nil {
		log.Printf("WARN: error preparing database statement to remove bookmark: %d: %v:\n", id, err)
		return
	}

	if _, err = stmt.Exec(id); err != nil {
		log.Printf("WARN: error executing database statement to remove bookmark: %d: %v:\n", id, err)
		return
	}

	go removeThumbLocally(id)
}

func TotalPageCount(database *sql.DB) int {
	rows, err := database.Query("SELECT COUNT(*) FROM bookmarks;")
	if err != nil {
		log.Println("WARN: could not get total count from database:", err)
		return -1
	}

	var rowCount int
	for rows.Next() {
		rows.Scan(&rowCount)
	}

	pageCount := rowCount / constants.PAGE_UPDATE_LIMIT
	if rowCount%constants.PAGE_UPDATE_LIMIT == 0 {
		pageCount -= 1
	}

	return pageCount
}

func ViewAllWebUI(database *sql.DB, pageNo int) []setup.Bookmark {
	var results []setup.Bookmark
	var result setup.Bookmark
	var modified time.Time

	stmt, err := database.Prepare("SELECT * FROM bookmarks ORDER BY id DESC LIMIT (?) OFFSET (?);")
	if err != nil {
		log.Println("WARN: error preparing database statement to get all bookmarks for web UI:", err)
		return results
	}

	pageOffset := pageNo * constants.PAGE_UPDATE_LIMIT

	rows, err := stmt.Query(constants.PAGE_UPDATE_LIMIT, pageOffset)
	if err != nil {
		log.Println("WARN: error executing database statement to get all bookmarks for web UI:", err)
		return results
	}

	for rows.Next() {
		result = setup.Bookmark{}
		rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.Category, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &modified)
		appendBookmarks(&results, result, modified)
	}
	return results
}

func ViewAll(database *sql.DB) []setup.Bookmark {
	var results []setup.Bookmark
	var result setup.Bookmark
	var modified time.Time

	rows, err := database.Query("SELECT * FROM bookmarks ORDER BY id DESC;")
	if err != nil {
		log.Println("WARN: error executing database statement to get all bookmarks:", err)
		return results
	}

	for rows.Next() {
		result = setup.Bookmark{}
		rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.Category, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &modified)
		appendBookmarks(&results, result, modified)
	}
	return results
}

func BackupViewAll(database *sql.DB) []setup.Bookmark {
	var results []setup.Bookmark
	var result setup.Bookmark
	var modified time.Time

	rows, err := database.Query("SELECT * FROM bookmarks ORDER BY id;")
	if err != nil {
		log.Println("WARN: error executing database statement to backup bookmarks:", err)
		return results
	}

	for rows.Next() {
		result = setup.Bookmark{}
		rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.Category, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &modified)
		appendBookmarks(&results, result, modified)
	}
	return results
}

func searchPageCount(database *sql.DB, query string, params []any) int {
	updatedQuery := strings.Replace(query, "*", "COUNT(*)", 1)
	params[len(params)-1] = 0

	var count int
	if err := database.QueryRow(updatedQuery, params...).Scan(&count); err != nil {
		log.Println("WARN: could not get page count of search query:", err)
	}

	pageCount := count / constants.PAGE_UPDATE_LIMIT
	if count%constants.PAGE_UPDATE_LIMIT == 0 {
		pageCount -= 1
	}

	return pageCount
}

func executeSearchQuery(database *sql.DB, searchType, searchTerm string, pageOffset int) (*sql.Rows, int, error) {
	var query string
	var params []any
	count := -1
	switch searchType {
	case "general":
		query = "SELECT * FROM bookmarks WHERE keywords LIKE (?) OR category LIKE (?) OR note LIKE (?) OR title LIKE (?) OR url LIKE (?) OR id LIKE (?) ORDER BY id DESC LIMIT (?) OFFSET (?);"
		params = []any{searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, constants.PAGE_UPDATE_LIMIT, pageOffset}
	case "hostname":
		query = "SELECT * FROM bookmarks WHERE url LIKE (?) ORDER BY id DESC LIMIT (?) OFFSET (?);"
		params = []any{searchTerm, constants.PAGE_UPDATE_LIMIT, pageOffset}
	case "keyword":
		query = "SELECT * FROM bookmarks WHERE keywords LIKE (?) ORDER BY id DESC LIMIT (?) OFFSET (?);"
		params = []any{searchTerm, constants.PAGE_UPDATE_LIMIT, pageOffset}
	case "category":
		query = "SELECT * FROM bookmarks WHERE category LIKE (?) ORDER BY id DESC LIMIT (?) OFFSET (?);"
		params = []any{searchTerm, constants.PAGE_UPDATE_LIMIT, pageOffset}
	default:
		return nil, count, fmt.Errorf("unrecognized search type: %s", searchType)
	}

	stmt, err := database.Prepare(query)
	if err != nil {
		return nil, count, err
	}

	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, count, err
	}

	count = searchPageCount(database, query, params)

	return rows, count, nil
}

func SearchFor(database *sql.DB, searchType, searchTerm string, pageNumber int) ([]setup.Bookmark, int) {
	var results []setup.Bookmark
	var result setup.Bookmark
	var modified time.Time

	if searchTerm == "" {
		return results, -1
	}
	searchTerm = "%" + searchTerm + "%"
	pageOffset := pageNumber * constants.PAGE_UPDATE_LIMIT

	rows, count, err := executeSearchQuery(database, searchType, searchTerm, pageOffset)
	if err != nil {
		log.Printf("WARN: error executing database query while searching for: %s: %v:\n", searchTerm, err)
		return results, -1
	}

	for rows.Next() {
		result = setup.Bookmark{}
		rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.Category, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &modified)
		appendBookmarks(&results, result, modified)
	}

	return results, count
}

func OpenSesame(database *sql.DB, searchTerm string) setup.Bookmark {
	var result setup.Bookmark
	var modified time.Time

	if searchTerm == "" {
		return result
	}
	searchTerm = "%" + searchTerm + "%"

	stmt, err := database.Prepare("SELECT * FROM bookmarks WHERE keywords LIKE (?) ORDER BY id DESC LIMIT (?);")
	if err != nil {
		log.Printf("WARN: failed to prepare database query to search keyword: %s: %v:\n", searchTerm, err)
		return result
	}

	row, err := stmt.Query(searchTerm, 1)
	if err != nil {
		log.Printf("WARN: failed to query database to search keyword: %s: %v:\n", searchTerm, err)
		return result
	}

	for row.Next() {
		row.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.Category, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &modified)
	}
	return result
}

func ViewSingleRow(database *sql.DB, id int64) (setup.Bookmark, error) {
	var rowResult setup.Bookmark
	var modified time.Time

	stmt, err := database.Prepare("SELECT * FROM bookmarks WHERE id=(?);")
	if err != nil {
		return rowResult, err
	}

	execResult, err := stmt.Query(id)
	if err != nil {
		return rowResult, err
	}

	for execResult.Next() {
		if err = execResult.Scan(&rowResult.ID, &rowResult.URL, &rowResult.Title, &rowResult.Note, &rowResult.Keywords, &rowResult.Category, &rowResult.Archived, &rowResult.SnapshotURL, &rowResult.ThumbURL, &modified); err != nil {
			return rowResult, err
		}
		rowResult.Modified = modified.Local().Format(constants.TIME_FORMAT)
	}

	if rowResult.URL == "" {
		return rowResult, fmt.Errorf("ID does not exist")
	}

	return rowResult, nil
}

// SearchByUrl function:
// Used in browser extension to check if the current tab is bookmarked
// or not, so appropriate extension icon is shown. Not related to SearchFor
// function.
func SearchByUrl(database *sql.DB, searchUrl string) (setup.Bookmark, error) {
	var urlResult setup.Bookmark

	stmt, err := database.Prepare("SELECT * FROM bookmarks WHERE url=(?);")
	if err != nil {
		return urlResult, err
	}

	execResult, err := stmt.Query(searchUrl)
	if err != nil {
		return urlResult, err
	}

	for execResult.Next() {
		if err = execResult.Scan(&urlResult.ID, &urlResult.URL, &urlResult.Title, &urlResult.Note, &urlResult.Keywords, &urlResult.Category, &urlResult.Archived, &urlResult.SnapshotURL, &urlResult.ThumbURL, &urlResult.Modified); err != nil {
			return urlResult, err
		}
	}

	return urlResult, nil
}

func GetAllCategories(database *sql.DB) ([]string, error) {
	var allCategories []string

	rows, err := database.Query("SELECT DISTINCT category FROM bookmarks ORDER BY id DESC;")
	if err != nil {
		return allCategories, err
	}

	var bkmCategory string
	for rows.Next() {
		rows.Scan(&bkmCategory)
		if bkmCategory != "" {
			allCategories = append(allCategories, bkmCategory)
		}
	}

	return allCategories, nil
}
