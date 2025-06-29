package db

import (
    "dalennod/internal/setup"
    "dalennod/internal/constants"
    "fmt"
    "time"
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
        ID:           info.ID,
        URL:          info.URL,
        Title:        info.Title,
        Note:         info.Note,
        Keywords:     info.Keywords,
        Category:     info.Category,
        Archived:     info.Archived,
        SnapshotURL:  info.SnapshotURL,
        ThumbURL:     info.ThumbURL,
        ByteThumbURL: info.ByteThumbURL,
        Modified:     modified.Local().Format(constants.TIME_FORMAT),
    })
}
