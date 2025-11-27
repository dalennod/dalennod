package db

import (
	"database/sql"
	"log"
	"math/rand/v2"

	"dalennod/internal/constants"
	"dalennod/internal/setup"
)

func recentsCount(database *sql.DB) int {
	rows, err := database.Query("SELECT COUNT(*) FROM recents")
	if err != nil {
		log.Println("WARN: querying database statement:", err)
		return -1
	}

	var count int
	for rows.Next() {
		rows.Scan(&count)
	}
	return count
}

func sanitizeRecents(database *sql.DB) {
	count := recentsCount(database)
	if count <= constants.RECENT_ENGAGE_LIMIT {
		log.Println("WARN: count less than RECENT_ENGAGE_LIMIT. no need to sanitize recents table")
		return
	}

	deleteLimit := count - constants.RECENT_ENGAGE_LIMIT
	stmt, err := database.Prepare("DELETE FROM recents WHERE id IN (SELECT id FROM recents ORDER BY lastAccessed ASC LIMIT (?));")
	if err != nil {
		log.Println("WARN: could not prepare database query:", err)
		return
	}
	if _, err := stmt.Exec(deleteLimit); err != nil {
		log.Println("WARN: could not execute query:", err)
		return
	}
}

func AddToRecents(database *sql.DB, bkmID int64) {
	// Clean up recents table of extra rows 5 out of 100 times AddToRecents is called.
	// This function will be called quite often, but it is not necessary to remove
	// extra rows on every call. That will have noticeable degradation in performance.
	if randNum := rand.IntN(100-1) + 1; randNum <= 5 {
		sanitizeRecents(database)
	}

	// Reference: https://sqlite.org/lang_upsert.html
	stmt, err := database.Prepare("INSERT INTO recents (bookmarkID, lastAccessed) VALUES (?, CURRENT_TIMESTAMP) ON CONFLICT(bookmarkID) DO UPDATE SET lastAccessed=CURRENT_TIMESTAMP;")
	if err != nil {
		log.Println("WARN: preparing database statement:", err)
		return
	}

	if _, err = stmt.Exec(bkmID); err != nil {
		log.Println("WARN: error executing database statement:", err)
		return
	}
}

func RecentInteractions(database *sql.DB) setup.RecentInteractions {
	var results setup.RecentInteractions

	stmt, err := database.Prepare("SELECT * FROM recents ORDER BY lastAccessed DESC LIMIT (?)")
	if err != nil {
		log.Println("WARN: could not prepare database query:", err)
		return results
	}

	rows, err := stmt.Query(constants.RECENT_ENGAGE_LIMIT)
	if err != nil {
		log.Println("WARN: could not execute query:", err)
		return results
	}

	for rows.Next() {
		rows.Scan(&results.ID, &results.BookmarkID, &results.LastAccessed)
		bkmRow, err := ViewSingleRow(database, results.BookmarkID)
		if err != nil {
			log.Printf("WARN: failed to get bookmark of ID %d: %v\n", results.BookmarkID, err)
			continue
		}
		results.Bookmarks = append(results.Bookmarks, bkmRow)
	}

	return results
}
