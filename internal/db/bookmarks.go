package db

import (
    "database/sql"
    "fmt"
    "strings"
    "time"

    "dalennod/internal/constants"
    "dalennod/internal/logger"
    "dalennod/internal/setup"
    "dalennod/internal/thumb_url"
)

func Add(database *sql.DB, bkmStruct setup.Bookmark) {
    if bkmStruct.ThumbURL == "" || len(bkmStruct.ByteThumbURL) == 0 {
        var err error
        bkmStruct.ThumbURL, bkmStruct.ByteThumbURL, err = thumb_url.GetPageThumb(bkmStruct.URL)
        if err != nil {
            logger.Error.Println("error getting webpage thumbnail. ERROR:", err)
        }
    }

    stmt, err := database.Prepare("INSERT INTO bookmarks (url, title, note, keywords, category, archived, snapshotURL, thumbURL, byteThumbURL) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);")
    if err != nil {
        logger.Error.Println("error preparing database statement. ERROR:", err)
        return
    }

    _, err = stmt.Exec(bkmStruct.URL, bkmStruct.Title, bkmStruct.Note, bkmStruct.Keywords, bkmStruct.Category, bkmStruct.Archived, bkmStruct.SnapshotURL, bkmStruct.ThumbURL, bkmStruct.ByteThumbURL)
    if err != nil {
        logger.Error.Println("error executing database statement. ERROR:", err)
        return
    }
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

    if bkmStruct.ThumbURL == "" || len(bkmStruct.ByteThumbURL) == 0 {
        var err error
        bkmStruct.ThumbURL, bkmStruct.ByteThumbURL, err = thumb_url.GetPageThumb(bkmStruct.URL)
        if err != nil {
            logger.Error.Println("error getting webpage thumbnail. ERROR:", err)
        }
    }

    stmt, err := database.Prepare("UPDATE bookmarks SET url=(?), title=(?), note=(?), keywords=(?), category=(?), archived=(?), snapshotURL=(?), thumbURL=(?), byteThumbURL=(?), modified=CURRENT_TIMESTAMP WHERE id=(?);")
    if err != nil {
        logger.Error.Println("error preparing database statement. ERROR:", err)
        return
    }

    _, err = stmt.Exec(bkmStruct.URL, bkmStruct.Title, bkmStruct.Note, bkmStruct.Keywords, bkmStruct.Category, bkmStruct.Archived, bkmStruct.SnapshotURL, bkmStruct.ThumbURL, bkmStruct.ByteThumbURL, bkmStruct.ID)
    if err != nil {
        logger.Error.Println("error executing database statement. ERROR:", err)
        return
    }
}

func updateCheck(database *sql.DB, bkmStruct setup.Bookmark) setup.Bookmark {
    oldBKMData, err := ViewSingleRow(database, bkmStruct.ID)
    if err != nil {
        logger.Error.Printf("error getting bookmark row data at ID %d. ERROR: %v\n", bkmStruct.ID, err)
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

func RefetchThumbnail(database *sql.DB, id int, thumbnail []byte) error {
    bkmStruct, err := ViewSingleRow(database, id)
    if err != nil {
        logger.Error.Printf("error getting bookmark row data at ID %d. ERROR: %v\n", id, err)
        return err
    }

    if thumbnail == nil {
        bkmStruct.ThumbURL, bkmStruct.ByteThumbURL, err = thumb_url.GetPageThumb(bkmStruct.URL)
        if err != nil || bkmStruct.ByteThumbURL == nil {
            logger.Error.Println("error getting webpage thumbnail. ERROR:", err)
            return err
        }
    } else {
        bkmStruct.ByteThumbURL = thumbnail
    }

    Update(database, bkmStruct, true)
    return nil
}

func Remove(database *sql.DB, id int) {
    stmt, err := database.Prepare("DELETE FROM bookmarks WHERE id=(?);")
    if err != nil {
        logger.Error.Println("error preparing database statement. ERROR:", err)
        return
    }

    if _, err = stmt.Exec(id); err != nil {
        logger.Error.Println("error executing database statement. ERROR:", err)
        return
    }
}

func TotalPageCount(database *sql.DB) int {
    rows, err := database.Query("SELECT COUNT(*) FROM bookmarks;")
    if err != nil {
        logger.Error.Println("error getting total count from database. ERROR:", err)
        return -1
    }

    var pageCount int
    for rows.Next() {
        rows.Scan(&pageCount)
    }
    pageCount = pageCount / constants.PAGE_UPDATE_LIMIT
    return pageCount
}

func ViewAllWebUI(database *sql.DB, pageNo int) []setup.Bookmark {
    var results []setup.Bookmark
    var result setup.Bookmark
    var modified time.Time

    stmt, err := database.Prepare("SELECT * FROM bookmarks ORDER BY id DESC LIMIT (?) OFFSET (?);")
    if err != nil {
        logger.Error.Println("error preparing database statement. ERROR:", err)
        return results
    }

    pageOffset := pageNo * constants.PAGE_UPDATE_LIMIT

    rows, err := stmt.Query(constants.PAGE_UPDATE_LIMIT, pageOffset)
    if err != nil {
        logger.Error.Println("error executing database statement. ERROR:", err)
        return results
    }

    for rows.Next() {
        result = setup.Bookmark{}
        rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.Category, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
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
        logger.Error.Println("error executing database statement. ERROR:", err)
        return results
    }

    for rows.Next() {
        result = setup.Bookmark{}
        rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.Category, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
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
        logger.Error.Println("error executing database statement. ERROR:", err)
        return results
    }

    for rows.Next() {
        result = setup.Bookmark{}
        rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.Category, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
        appendBookmarks(&results, result, modified)
    }
    return results
}

func searchPageCount(database *sql.DB, query string, params []interface{}) int {
    updatedQuery := strings.Replace(query, "*", "COUNT(*)", 1)
    params[len(params)-1] = 0
    var count int
    if err := database.QueryRow(updatedQuery, params...).Scan(&count); err != nil {
        logger.Error.Println("error getting page count of search query. ERROR:", err)
    }
    count = count / constants.PAGE_UPDATE_LIMIT
    return count
}

func executeSearchQuery(database *sql.DB, searchType, searchTerm string, pageOffset int) (*sql.Rows, int, error) {
    var query string
    var params []interface{}
    count := -1
    switch searchType {
    case "general":
        query = "SELECT * FROM bookmarks WHERE keywords LIKE (?) OR category LIKE (?) OR note LIKE (?) OR title LIKE (?) OR url LIKE (?) ORDER BY id DESC LIMIT (?) OFFSET (?);"
        params = []interface{}{searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, constants.PAGE_UPDATE_LIMIT, pageOffset}
    case "hostname":
        query = "SELECT * FROM bookmarks WHERE url LIKE (?) ORDER BY id DESC LIMIT (?) OFFSET (?);"
        params = []interface{}{searchTerm, constants.PAGE_UPDATE_LIMIT, pageOffset}
    case "keyword":
        query = "SELECT * FROM bookmarks WHERE keywords LIKE (?) ORDER BY id DESC LIMIT (?) OFFSET (?);"
        params = []interface{}{searchTerm, constants.PAGE_UPDATE_LIMIT, pageOffset}
    case "category":
        query = "SELECT * FROM bookmarks WHERE category LIKE (?) ORDER BY id DESC LIMIT (?) OFFSET (?);"
        params = []interface{}{searchTerm, constants.PAGE_UPDATE_LIMIT, pageOffset}
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
        logger.Error.Println("error executing database query. ERROR:", err)
        return results, -1
    }

    for rows.Next() {
        result = setup.Bookmark{}
        rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.Category, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
        appendBookmarks(&results, result, modified)
    }

    return results, count
}

func ViewSingleRow(database *sql.DB, id int) (setup.Bookmark, error) {
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
        if err = execRes.Scan(&rowResult.ID, &rowResult.URL, &rowResult.Title, &rowResult.Note, &rowResult.Keywords, &rowResult.Category, &rowResult.Archived, &rowResult.SnapshotURL, &rowResult.ThumbURL, &rowResult.ByteThumbURL, &modified); err != nil {
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

    execRes, err := stmt.Query(searchUrl)
    if err != nil {
        return urlResult, err
    }

    for execRes.Next() {
        if err = execRes.Scan(&urlResult.ID, &urlResult.URL, &urlResult.Title, &urlResult.Note, &urlResult.Keywords, &urlResult.Category, &urlResult.Archived, &urlResult.SnapshotURL, &urlResult.ThumbURL, &urlResult.ByteThumbURL, &urlResult.Modified); err != nil {
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
