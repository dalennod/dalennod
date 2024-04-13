package user_input

import (
	"bufio"
	"dalennod/internal/archive"
	"dalennod/internal/backup"
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"dalennod/internal/server"
	"dalennod/internal/setup"
	"dalennod/internal/thumb_url"
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

var (
	database *sql.DB
	flagVals setup.FlagValues
)

func UserInput(data *sql.DB) {
	enableLogs()

	database = data

	flagVals = setup.ParseFlags()

	switch true {
	case flagVals.ViewAll:
		db.ViewAll(database, false)
	case flagVals.StartServer:
		server.Start(database)
	case flagVals.AddEntry:
		addInput("", "", "", "", "", "", false, 0)
	case flagVals.Backup && flagVals.JSONOut:
		backup.JSONOut(database)
	case flagVals.Backup:
		backup.GDrive()
	}

	if flagVals.RemoveID != "" {
		removeInput(flagVals.RemoveID)
	} else if flagVals.UpdateID != "" {
		updateInput(flagVals.UpdateID)
	} else if flagVals.ViewID != "" {
		viewInput(flagVals.ViewID)
	}
}

func addInput(url, title, note, keywords, group, archived string, update bool, id int) {
	var (
		archiveResult bool   = false
		snapshotURL   string = ""
		scanner              = bufio.NewScanner(os.Stdin)
	)

	fmt.Print("URL to save: ")
	scanner.Scan()
	url = scanner.Text()

	thumbURL, err := thumb_url.GetPageThumb(url)
	if err != nil {
		thumbURL = url
	}

	fmt.Print("Title for the bookmark: ")
	scanner.Scan()
	title = scanner.Text()

	fmt.Print("Notes/log reason for bookmark: ")
	scanner.Scan()
	note = scanner.Text()

	fmt.Print("Keywords for searching later: ")
	scanner.Scan()
	keywords = scanner.Text()

	fmt.Print("Group to store the bookmark into: ")
	scanner.Scan()
	group = scanner.Text()

	fmt.Print("Archive URL? (y/N): ")
	scanner.Scan()
	archived = scanner.Text()

	if !update {
		switch archived {
		case "y", "Y":
			archiveResult, snapshotURL = archive.SendSnapshot(url)
			if archiveResult {
				db.Add(database, url, title, note, keywords, group, true, snapshotURL, thumbURL)
			} else {
				logger.Warn.Println("Snapshot failed.")
				db.Add(database, url, title, note, keywords, group, false, snapshotURL, thumbURL)
			}
		case "n", "N":
			db.Add(database, url, title, note, keywords, group, false, snapshotURL, thumbURL)
		default:
			db.Add(database, url, title, note, keywords, group, false, snapshotURL, thumbURL)
			logger.Warn.Println("Invalid input for archive request. URL has not been archived.")
		}
	} else {
		switch archived {
		case "y", "Y":
			archiveResult, snapshotURL = archive.SendSnapshot(url)
			if archiveResult {
				db.Update(database, url, title, note, keywords, group, id, true, false, snapshotURL)
			} else {
				logger.Warn.Println("Snapshot failed.")
				db.Update(database, url, title, note, keywords, group, id, false, false, snapshotURL)
			}
		case "n", "N":
			db.Update(database, url, title, note, keywords, group, id, false, false, snapshotURL)
		default:
			db.Update(database, url, title, note, keywords, group, id, false, false, snapshotURL)
			logger.Warn.Println("Invalid input for archive request. URL has not been archived.")
		}
	}
}

func updateInput(updateID string) {
	var (
		id, url, title, note, keywords, group, archived, confirm string
		scanner                                                  = bufio.NewScanner(os.Stdin)
	)

	// fmt.Print("ID of bookmark to update: ")
	// scanner.Scan()
	// id = scanner.Text()

	id = updateID

	idToINT, err := strconv.Atoi(id)
	if err != nil {
		logger.Error.Println("invalid input")
	}

	db.ViewSingleRow(database, idToINT, false)

	fmt.Print("Update this entry? (y/N): ")
	scanner.Scan()
	confirm = scanner.Text()

	switch confirm {
	case "y", "Y":
		fmt.Println("Leave empty to retain old information.")
		addInput(url, title, note, keywords, group, archived, true, idToINT)
	case "n", "N":
		return
	default:
		logger.Info.Println("Invalid input received:", confirm)
		fmt.Println("Invalid input. Exiting.")
		return
	}
}

func removeInput(removeID string) {
	var (
		id, confirm string
		scanner     = bufio.NewScanner(os.Stdin)
	)

	// fmt.Print("ID to remove: ")
	// scanner.Scan()
	// id = scanner.Text()

	id = removeID

	idToINT, err := strconv.Atoi(id)
	if err != nil {
		logger.Error.Println("Invalid input.")
	}

	db.ViewSingleRow(database, idToINT, false)

	fmt.Print("Remove this entry? (y/n): ")
	scanner.Scan()
	confirm = scanner.Text()

	switch confirm {
	case "y", "Y":
		db.Remove(database, idToINT)
	case "n", "N":
		return
	default:
		logger.Info.Println("Invalid input received:", confirm)
		fmt.Println("Invalid input. Exiting.")
		return
	}
}

func viewInput(viewID string) {
	idToINT, err := strconv.Atoi(viewID)
	if err != nil {
		logger.Error.Println("Invalid input.")
	}
	db.ViewSingleRow(database, idToINT, false)
}

func enableLogs() {
	logger.Enable()
	cfgDir, _ := setup.ConfigDir()
	logDir, _ := setup.CacheDir()
	logger.Info.Printf("Database and config directory: %s\n", cfgDir)
	logger.Info.Printf("Error logs directory: %s\n", logDir)
}
