package db

import (
	"bufio"
	"dalennod/internal/archive"
	"dalennod/internal/logger"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var db *sql.DB

func UserInput(database *sql.DB) {
	db = database

	var input string
	fmt.Print("Add, update, remove or view entries? : ")
	fmt.Scanln(&input)

	switch strings.ToLower(input) {
	case "add", "a":
		addInput("", "", "", "", "", "", false, 0)
	case "update", "u":
		updateInput()
	case "remove", "r":
		removeInput()
	case "view", "v":
		viewInput()
	case "q":
		os.Exit(0)
	default:
		logger.Warn.Println("invalid input")
	}
}

func addInput(url, title, note, keywords, group, archived string, update bool, id int) {
	var (
		archiveResult bool   = false
		snapshotURL   string = ""
		scanner              = bufio.NewScanner(os.Stdin)
	)

	logger.Info.Println(url, title, note, keywords, group, archived)

	fmt.Print("URL to save: ")
	scanner.Scan()
	url = scanner.Text()

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
				Add(db, url, title, note, keywords, group, true, snapshotURL)
			} else {
				logger.Warn.Println("Snapshot failed.")
				Add(db, url, title, note, keywords, group, false, snapshotURL)
			}
		case "n", "N":
			Add(db, url, title, note, keywords, group, false, snapshotURL)
		default:
			Add(db, url, title, note, keywords, group, false, snapshotURL)
			logger.Warn.Println("Invalid input for archive request. URL has not been archived.")
		}
	} else {
		switch archived {
		case "y", "Y":
			archiveResult, snapshotURL = archive.SendSnapshot(url)
			if archiveResult {
				Update(db, url, title, note, keywords, group, id, true, false, snapshotURL)
			} else {
				logger.Warn.Println("Snapshot failed.")
				Update(db, url, title, note, keywords, group, id, false, false, snapshotURL)
			}
		case "n", "N":
			Update(db, url, title, note, keywords, group, id, false, false, snapshotURL)
		default:
			Update(db, url, title, note, keywords, group, id, false, false, snapshotURL)
			logger.Warn.Println("Invalid input for archive request. URL has not been archived.")
		}
	}
}

func updateInput() {
	var (
		id, url, title, note, keywords, group, archived, confirm string
		scanner                                                  = bufio.NewScanner(os.Stdin)
	)

	fmt.Print("ID of bookmark to update: ")
	scanner.Scan()
	id = scanner.Text()

	idToINT, err := strconv.Atoi(id)
	if err != nil {
		logger.Error.Println("invalid input")
	}

	ViewSingleRow(db, idToINT, false)

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

func removeInput() {
	var (
		id, confirm string
		scanner     = bufio.NewScanner(os.Stdin)
	)

	fmt.Print("ID to remove: ")
	scanner.Scan()
	id = scanner.Text()

	idToINT, err := strconv.Atoi(id)
	if err != nil {
		logger.Error.Println("Invalid input.")
	}

	ViewSingleRow(db, idToINT, false)

	fmt.Print("Remove this entry? (y/n): ")
	scanner.Scan()
	confirm = scanner.Text()

	switch confirm {
	case "y", "Y":
		Remove(db, idToINT)
	case "n", "N":
		return
	default:
		logger.Info.Println("Invalid input received:", confirm)
		fmt.Println("Invalid input. Exiting.")
		return
	}
}

func viewInput() {
	var (
		input   string
		scanner = bufio.NewScanner(os.Stdin)
	)

	fmt.Print("Enter specific ID or All: ")
	scanner.Scan()
	input = scanner.Text()

	switch strings.ToLower(input) {
	case "all", "a":
		ViewAll(db, false)
	default:
		idToINT, err := strconv.Atoi(input)
		if err != nil {
			logger.Error.Println("Invalid input.")
		}
		ViewSingleRow(db, idToINT, false)
	}
}
