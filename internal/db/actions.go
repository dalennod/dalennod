package db

import (
    "dalennod/internal/logger"
    "dalennod/internal/setup"
    "dalennod/internal/thumb_url"
    "dalennod/internal/constants"
    "database/sql"
    "fmt"
    "time"
)

func Add(database *sql.DB, bmStruct setup.Bookmark) {
    if bmStruct.ThumbURL == "" || len(bmStruct.ByteThumbURL) == 0 {
        var err error
        bmStruct.ThumbURL, bmStruct.ByteThumbURL, err = thumb_url.GetPageThumb(bmStruct.URL)
        if err != nil {
            logger.Error.Println("error getting webpage thumbnail. ERROR:", err)
        }
    }

    stmt, err := database.Prepare("INSERT INTO bookmarks (url, title, note, keywords, bmGroup, archived, snapshotURL, thumbURL, byteThumbURL) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);")
    if err != nil {
        logger.Error.Println("error preparing database statement. ERROR:", err)
        return
    }

    _, err = stmt.Exec(bmStruct.URL, bmStruct.Title, bmStruct.Note, bmStruct.Keywords, bmStruct.BmGroup, bmStruct.Archived, bmStruct.SnapshotURL, bmStruct.ThumbURL, bmStruct.ByteThumbURL)
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
func Update(database *sql.DB, bmStruct setup.Bookmark, serverCall bool) {
    if !serverCall {
        bmStruct = updateCheck(database, bmStruct)
    }

    if bmStruct.ThumbURL == "" || len(bmStruct.ByteThumbURL) == 0 {
        var err error
        bmStruct.ThumbURL, bmStruct.ByteThumbURL, err = thumb_url.GetPageThumb(bmStruct.URL)
        if err != nil {
            logger.Error.Println("error getting webpage thumbnail. ERROR:", err)
        }
    }

    stmt, err := database.Prepare("UPDATE bookmarks SET url=(?), title=(?), note=(?), keywords=(?), bmGroup=(?), archived=(?), snapshotURL=(?), thumbURL=(?), byteThumbURL=(?), modified=(?) WHERE id=(?);")
    if err != nil {
        logger.Error.Println("error preparing database statement. ERROR:", err)
        return
    }

    _, err = stmt.Exec(bmStruct.URL, bmStruct.Title, bmStruct.Note, bmStruct.Keywords, bmStruct.BmGroup, bmStruct.Archived, bmStruct.SnapshotURL, bmStruct.ThumbURL, bmStruct.ByteThumbURL, time.Now().UTC().Format(constants.TIME_FORMAT), bmStruct.ID)
    if err != nil {
        logger.Error.Println("error executing database statement. ERROR:", err)
        return
    }
}

func updateCheck(database *sql.DB, bmStruct setup.Bookmark) setup.Bookmark {
    oldBmData, err := ViewSingleRow(database, bmStruct.ID)
    if err != nil {
        logger.Error.Printf("error getting bookmark row data at ID %d. ERROR: %v\n", bmStruct.ID, err)
        return bmStruct
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
    bmStruct, err := ViewSingleRow(database, id)
    if err != nil {
        logger.Error.Printf("error getting bookmark row data at ID %d. ERROR: %v\n", id, err)
        return err
    }

    if thumbnail == nil {
        bmStruct.ThumbURL, bmStruct.ByteThumbURL, err = thumb_url.GetPageThumb(bmStruct.URL)
        if err != nil || bmStruct.ByteThumbURL == nil {
            logger.Error.Println("error getting webpage thumbnail. ERROR:", err)
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
        logger.Error.Println("error preparing database statement. ERROR:", err)
        return
    }

    _, err = stmt.Exec(id)
    if err != nil {
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
        rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
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
        rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
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
        rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
        appendBookmarks(&results, result, modified)
    }
    return results
}

func executeSearchQuery(database *sql.DB, searchType, searchTerm string, pageOffset int) (*sql.Rows, error) {
    var query string;
    var params []interface{};
    switch searchType {
    case "general":
        query = "SELECT * FROM bookmarks WHERE keywords LIKE (?) OR bmGroup LIKE (?) OR note LIKE (?) OR title LIKE (?) OR url LIKE (?) ORDER BY id DESC LIMIT (?) OFFSET (?);";
        params = []interface{}{searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, constants.PAGE_UPDATE_LIMIT, pageOffset};
    case "hostname":
        query = "SELECT * FROM bookmarks WHERE url LIKE (?) ORDER BY id DESC LIMIT (?) OFFSET (?);";
        params = []interface{}{searchTerm, constants.PAGE_UPDATE_LIMIT, pageOffset};
    case "keyword":
        query = "SELECT * FROM bookmarks WHERE keywords LIKE (?) ORDER BY id DESC LIMIT (?) OFFSET (?);";
        params = []interface{}{searchTerm, constants.PAGE_UPDATE_LIMIT, pageOffset};
    case "group":
        query = "SELECT * FROM bookmarks WHERE bmGroup LIKE (?) ORDER BY id DESC LIMIT (?) OFFSET (?);";
        params = []interface{}{searchTerm, constants.PAGE_UPDATE_LIMIT, pageOffset};
    default:
        return nil, fmt.Errorf("unrecognized search type: %s", searchType);
    }

    stmt, err := database.Prepare(query);
    if err != nil {
        return nil, err;
    }

    rows, err := stmt.Query(params...);
    if err != nil {
        return nil, err;
    }

    return rows, nil;
}

func SearchFor(database *sql.DB, searchType, searchTerm string, pageNumber int) []setup.Bookmark {
    var results []setup.Bookmark;
    var result setup.Bookmark;
    var modified time.Time;

    if searchTerm == "" {
        return results;
    }
    searchTerm = "%" + searchTerm + "%";
    pageOffset := pageNumber * constants.PAGE_UPDATE_LIMIT;

    rows, err := executeSearchQuery(database, searchType, searchTerm, pageOffset)
    if err != nil {
        logger.Error.Println("error executing database query. ERROR:", err)
        return results
    }

    for rows.Next() {
        result = setup.Bookmark{}
        rows.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
        appendBookmarks(&results, result, modified)
    }

    return results;
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
        if err = execRes.Scan(&rowResult.ID, &rowResult.URL, &rowResult.Title, &rowResult.Note, &rowResult.Keywords, &rowResult.BmGroup, &rowResult.Archived, &rowResult.SnapshotURL, &rowResult.ThumbURL, &rowResult.ByteThumbURL, &modified); err != nil {
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
        if err = execRes.Scan(&urlResult.ID, &urlResult.URL, &urlResult.Title, &urlResult.Note, &urlResult.Keywords, &urlResult.BmGroup, &urlResult.Archived, &urlResult.SnapshotURL, &urlResult.ThumbURL, &urlResult.ByteThumbURL, &urlResult.Modified); err != nil {
            return urlResult, err;
        }
    }

    return urlResult, nil
}

func GetAllGroups(database *sql.DB) ([]string, error) {
    var allGroups []string

    rows, err := database.Query("SELECT DISTINCT bmGroup FROM bookmarks ORDER BY id DESC;")
    if err != nil {
        return allGroups, err
    }

    var bmGroup string
    for rows.Next() {
        rows.Scan(&bmGroup)
        if bmGroup != "" {
            allGroups = append(allGroups, bmGroup)
        }
    }

    return allGroups, nil
}
