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

func Update(database *sql.DB, bmStruct setup.Bookmark, serverCall bool) {
    if !serverCall {
        bmStruct = updateCheck(database, bmStruct)
    }

    stmt, err := database.Prepare("UPDATE bookmarks SET url=(?), title=(?), note=(?), keywords=(?), bmGroup=(?), archived=(?), snapshotURL=(?), thumbURL=(?), byteThumbURL=(?) WHERE id=(?);")
    if err != nil {
        logger.Error.Println("error preparing database statement. ERROR:", err)
        return
    }

    _, err = stmt.Exec(bmStruct.URL, bmStruct.Title, bmStruct.Note, bmStruct.Keywords, bmStruct.BmGroup, bmStruct.Archived, bmStruct.SnapshotURL, bmStruct.ThumbURL, bmStruct.ByteThumbURL, bmStruct.ID)
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
        logger.Error.Println("error preparing database statement. ERROR:", err)
        return results
    }

    execRes, err := stmt.Query(keyword, keyword, keyword, keyword, keyword)
    if err != nil {
        logger.Error.Println("error executing database statement. ERROR:", err)
        return results
    }

    for execRes.Next() {
        result = setup.Bookmark{}
        execRes.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
        appendBookmarks(&results, result, modified)
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
        logger.Error.Println("error preparing database statement. ERROR:", err)
        return results
    }

    execRes, err := stmt.Query(keyword)
    if err != nil {
        logger.Error.Println("error executing database statement. ERROR:", err)
        return results
    }

    for execRes.Next() {
        result = setup.Bookmark{}
        execRes.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
        appendBookmarks(&results, result, modified)
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
        logger.Error.Println("error preparing database statement. ERROR:", err)
        return results
    }

    execRes, err := stmt.Query(keyword)
    if err != nil {
        logger.Error.Println("error executing database statement. ERROR:", err)
        return results
    }

    for execRes.Next() {
        result = setup.Bookmark{}
        execRes.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
        appendBookmarks(&results, result, modified)
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
        logger.Error.Println("error preparing database statement. ERROR:", err)
        return results
    }

    execRes, err := stmt.Query(hostname)
    if err != nil {
        logger.Error.Println("error executing database statement. ERROR:", err)
        return results
    }

    for execRes.Next() {
        result = setup.Bookmark{}
        execRes.Scan(&result.ID, &result.URL, &result.Title, &result.Note, &result.Keywords, &result.BmGroup, &result.Archived, &result.SnapshotURL, &result.ThumbURL, &result.ByteThumbURL, &modified)
        appendBookmarks(&results, result, modified)
    }

    return results
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
        if err = execRes.Scan(
            &rowResult.ID,
            &rowResult.URL,
            &rowResult.Title,
            &rowResult.Note,
            &rowResult.Keywords,
            &rowResult.BmGroup,
            &rowResult.Archived,
            &rowResult.SnapshotURL,
            &rowResult.ThumbURL,
            &rowResult.ByteThumbURL,
            &modified,
        ); err != nil {
            return rowResult, err
        }
        rowResult.Modified = modified.Local().Format(constants.TIME_FORMAT)
    }

    if rowResult.URL == "" {
        return rowResult, fmt.Errorf("ID does not exist")
    }

    return rowResult, nil
}

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
        if err = execRes.Scan(
            &urlResult.ID,
            &urlResult.URL,
            &urlResult.Title,
            &urlResult.Note,
            &urlResult.Keywords,
            &urlResult.BmGroup,
            &urlResult.Archived,
            &urlResult.SnapshotURL,
            &urlResult.ThumbURL,
            &urlResult.ByteThumbURL,
            &urlResult.Modified,
        ); err != nil {
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
