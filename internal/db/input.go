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
	fmt.Print("Add, remove or view entries? : ")
	fmt.Scanln(&input)

	switch strings.ToLower(input) {
	case "add", "a":
		addInput()
	case "remove", "r":
		removeInput()
	case "view", "v":
		viewInput()
	default:
		logger.Warn.Println("invalid input")
	}
}

func addInput() {
	var (
		url, title, note, keywords, group, archived string
		scanner                                     = bufio.NewScanner(os.Stdin)
	)

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

	switch archived {
	case "y", "Y":
		Add(db, url, title, note, keywords, group, 1)
		archive.SendSnapshot(url) // need error checking
	case "n", "N":
		Add(db, url, title, note, keywords, group, 0)
	default:
		Add(db, url, title, note, keywords, group, 0)
		logger.Warn.Println("Invalid input for archive request. URL has not been archived.")
	}

}

func removeInput() {
	var (
		id      string
		scanner = bufio.NewScanner(os.Stdin)
	)

	fmt.Print("ID to remove: ")
	scanner.Scan()
	id = scanner.Text()

	idToINT, err := strconv.Atoi(id)
	if err != nil {
		logger.Error.Println("Invalid input.")
	}

	// double check with user using ViewSingleRow here
	Remove(db, idToINT)
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
		_ = ViewAll(db, "c")
	default:
		idToINT, err := strconv.Atoi(input)
		if err != nil {
			logger.Error.Println("Invalid input.")
		}
		ViewSingleRow(db, idToINT)
	}
}
