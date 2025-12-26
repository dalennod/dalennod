package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"dalennod/internal/constants"
	"dalennod/internal/setup"
	"dalennod/internal/thumb_url"
)

func PrintRow(bookmarkRow setup.Bookmark) {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("  #%d | %s\n", bookmarkRow.ID, bookmarkRow.Modified))
	sb.WriteString(fmt.Sprintf("Title       : %s\n", bookmarkRow.Title))
	sb.WriteString(fmt.Sprintf("URL         : %s\n", bookmarkRow.URL))
	sb.WriteString(fmt.Sprintf("Note        : %s\n", bookmarkRow.Note))
	sb.WriteString(fmt.Sprintf("Keywords    : %s\n", bookmarkRow.Keywords))
	sb.WriteString(fmt.Sprintf("Category    : %s\n", bookmarkRow.Category))
	sb.WriteString(fmt.Sprintf("Archived?   : %t\n", bookmarkRow.Archived))
	if bookmarkRow.Archived {
		sb.WriteString(fmt.Sprintf("Archive URL : %s\n", bookmarkRow.SnapshotURL))
	}

	fmt.Println(sb.String())
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

func saveThumbLocally(id int64, thumbURL string) {
	if thumbURL != "" {
		err := thumb_url.DownThumb(id, thumbURL)
		if err != nil {
			log.Println("WARN: could not save thumbnail locally:", err)
			return
		}
	}
}

func removeThumbLocally(id int64) {
	err := os.Remove(filepath.Join(constants.THUMBNAILS_PATH, strconv.FormatInt(id, 10)))
	if err != nil {
		log.Printf("WARN: could not remove thumbnail locally from for bookmark ID %d: %v\n", id, err)
	}
}
