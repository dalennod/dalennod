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
	fmt.Printf("  #%d | %s\n", bookmarkRow.ID, bookmarkRow.Modified)
	fmt.Println("Title       : ", bookmarkRow.Title)
	fmt.Println("URL         : ", bookmarkRow.URL)
	if strings.Contains(bookmarkRow.Note, "\n") {
		bkmNoteNewlineSplit := strings.Split(bookmarkRow.Note, "\n")
		for i, n := range bkmNoteNewlineSplit {
			if i == 0 {
				fmt.Println("Note        : ", n)
				continue
			}
			if n == "\n" || n == "" {
				continue
			}
			fmt.Println("            : ", n)
		}
	} else {
		fmt.Println("Note        : ", bookmarkRow.Note)
	}
	fmt.Println("Keywords    : ", bookmarkRow.Keywords)
	fmt.Println("Category    : ", bookmarkRow.Category)
	fmt.Println("Archived?   : ", bookmarkRow.Archived)
	if bookmarkRow.Archived {
		fmt.Println("Archive URL : ", bookmarkRow.SnapshotURL)
	}
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
