package db

import (
    "database/sql"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "time"

    "dalennod/internal/constants"
    "dalennod/internal/logger"
    "dalennod/internal/setup"
    "dalennod/internal/thumb_url"
)

func PrintRow(bookmarkRow setup.Bookmark) {
    fmt.Printf("    #%d -- %s\nTitle:\t\t%s\nURL:\t\t%s\nNote:\t\t%s\nKeywords:\t%s\nCategory:\t\t%s\nArchived?:\t%t\nArchive URL:\t%s\n\n",
        bookmarkRow.ID,
        bookmarkRow.Modified,
        bookmarkRow.Title,
        bookmarkRow.URL,
        bookmarkRow.Note,
        bookmarkRow.Keywords,
        bookmarkRow.Category,
        bookmarkRow.Archived,
        bookmarkRow.SnapshotURL)
}

func appendBookmarks(b *[]setup.Bookmark, info setup.Bookmark, modified time.Time) {
    *b = append(*b, setup.Bookmark{
        ID:          info.ID,
        URL:         info.URL,
        Title:       info.Title,
        Note:        info.Note,
        Keywords:    info.Keywords,
        Category:    info.Category,
        Archived:    info.Archived,
        SnapshotURL: info.SnapshotURL,
        ThumbURL:    info.ThumbURL,
        Modified:    modified.Local().Format(constants.TIME_FORMAT),
    })
}

func saveThumbLocally(execResult sql.Result, thumbURL string) {
    id, err := execResult.LastInsertId()
    if err != nil {
        logger.Warn.Println("could not get last insert ID. ERROR:", err)
        return
    }

    if thumbURL != "" {
        err := thumb_url.DownThumb(id, thumbURL)
        if err != nil {
            logger.Warn.Println("could not save thumbnail locally. ERROR:", err)
            return
        }
    }
}

func removeThumbLocally(id int) {
    err := os.Remove(filepath.Join(constants.THUMBNAILS_PATH, strconv.Itoa(id)))
    if err != nil {
        logger.Warn.Printf("could not remove thumbnail locally from %s for ID %d. ERROR: %v", constants.THUMBNAILS_PATH, id, err)
    }
}
