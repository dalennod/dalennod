package db

import (
	"database/sql"
	"math/rand/v2"

	"dalennod/internal/constants"
	"dalennod/internal/logger"
	"dalennod/internal/setup"
)

func recentsCount(database *sql.DB) int {
	rows, err := database.Query("SELECT COUNT(*) FROM recents")
	if err != nil {
		logger.Error.Println("error preparing database statement. ERROR:", err)
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
		logger.Info.Println("count less than RECENT_ENGAGE_LIMIT")
		return
	}

	deleteLimit := count - constants.RECENT_ENGAGE_LIMIT
	stmt, err := database.Prepare("DELETE FROM recents WHERE id IN (SELECT id FROM recents ORDER BY lastAccessed ASC LIMIT (?));")
	if err != nil {
		logger.Error.Println("could not prepare database query. ERROR:", err)
		return
	}
	if _, err := stmt.Exec(deleteLimit); err != nil {
		logger.Error.Println("could not execute query. ERROR:", err)
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
		logger.Error.Println("error preparing database statement. ERROR:", err)
		return
	}

	if _, err = stmt.Exec(bkmID); err != nil {
		logger.Error.Println("error executing database statement. ERROR:", err)
		return
	}
}

func RecentInteractions(database *sql.DB) setup.RecentInteractions {
	var results setup.RecentInteractions

	stmt, err := database.Prepare("SELECT * FROM recents ORDER BY lastAccessed DESC LIMIT (?)")
	if err != nil {
		logger.Error.Println("could not prepare database query. ERROR:", err)
		return results
	}

	rows, err := stmt.Query(constants.RECENT_ENGAGE_LIMIT)
	if err != nil {
		logger.Error.Println("could not execute query. ERROR:", err)
		return results
	}

	for rows.Next() {
		rows.Scan(&results.ID, &results.BookmarkID, &results.LastAccessed)
		bkmRow, err := ViewSingleRow(database, results.BookmarkID)
		if err != nil {
			logger.Warn.Printf("failed to get bookmark of ID %d: %v\n", results.BookmarkID, err)
			continue
		}
		results.Bookmarks = append(results.Bookmarks, bkmRow)
	}

	return results
}
