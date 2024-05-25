package user_input

import (
	"bufio"
	"bytes"
	"dalennod/internal/archive"
	"dalennod/internal/backup"
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"dalennod/internal/server"
	"dalennod/internal/setup"
	"dalennod/internal/thumb_url"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"

	"golang.org/x/net/html"
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
	}

	if flagVals.RemoveID != "" {
		removeInput(flagVals.RemoveID)
	} else if flagVals.UpdateID != "" {
		updateInput(flagVals.UpdateID)
	} else if flagVals.ViewID != "" {
		viewInput(flagVals.ViewID)
	} else if flagVals.Import != "" && flagVals.Firefox {
		importFirefoxInput(flagVals.Import)
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

	url, title, note, keywords, group, _ = db.ViewSingleRow(database, idToINT, false)
	if url == "" {
		fmt.Println("id does not exist")
		logger.Info.Println("id does not exist")
		return
	}

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

	url, _, _, _, _, _ := db.ViewSingleRow(database, idToINT, false)
	if url == "" {
		fmt.Println("id does not exist")
		logger.Info.Println("id does not exist")
		return
	}

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

// try to figure out how to import in Group values too at some point
func importFirefoxInput(file string) {
	rf, err := os.ReadFile(file)
	if err != nil {
		logger.Error.Printf("couldn't open file: [error: %v]", err)
	}

	var parsedBookmarks []setup.Bookmark
	parsedBookmarks, err = parseFfInputFile(bytes.NewReader(rf), parsedBookmarks)
	if err != nil {
		logger.Error.Fatalln("parsing error:", err)
	}
	var parsedBookmarksCount = len(parsedBookmarks)

	for i, parsedBookmark := range parsedBookmarks {
		db.Add(database, parsedBookmark.URL, parsedBookmark.Title, "", parsedBookmark.Keywords, "", false, "", parsedBookmark.ThumbURL)
		fmt.Printf("\rAdded %d / %d", i+1, parsedBookmarksCount)
	}
	fmt.Println()
}

func parseFfInputFile(htmlImport io.Reader, parsedBookmarks []setup.Bookmark) ([]setup.Bookmark, error) {
	parseHtmlImport, err := html.Parse(htmlImport)
	if err != nil {
		return parsedBookmarks, err
	}
	var processHtmlImport func(n *html.Node)
	processHtmlImport = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			parsedBookmarks = processNode(n, parsedBookmarks)
		}
		for childNode := n.FirstChild; childNode != nil; childNode = childNode.NextSibling {
			processHtmlImport(childNode)
		}
	}
	processHtmlImport(parseHtmlImport)

	return parsedBookmarks, nil
}

func processNode(n *html.Node, parsedBookmarks []setup.Bookmark) []setup.Bookmark {
	// var url, thumbUrl, addDate, tags, keywords, title string
	var url, thumbUrl, tags, keywords, title string

	switch n.Data {
	case "a":
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				url = attr.Val
			}
			if attr.Key == "icon_uri" {
				thumbUrl = attr.Val
			}
			// if attr.Key == "add_date" {
			// 	addDate = attr.Val
			// }
			if attr.Key == "tags" {
				tags = attr.Val
			}
			if attr.Key == "shortcuturl" {
				keywords = attr.Val
			}
		}
	}
	for childNode := n.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		title = childNode.Data

		// addDateInt, err := strconv.ParseInt(addDate, 10, 64)
		// if err != nil {
		// 	logger.Error.Println("error:", err)
		// 	addDateInt = 1
		// }

		parsedBookmarks = append(parsedBookmarks, setup.Bookmark{
			URL:      url,
			ThumbURL: thumbUrl,
			// Modified: time.Unix(addDateInt, 0).Local().Format("2006-01-02 15:04:05"),
			Keywords: tags + keywords,
			Title:    title,
		})

		processNode(childNode, parsedBookmarks)
	}

	return parsedBookmarks
}

func enableLogs() {
	logger.Enable()
	cfgDir, _ := setup.ConfigDir()
	logDir, _ := setup.CacheDir()
	logger.Info.Printf("Database and config directory: %s\n", cfgDir)
	logger.Info.Printf("Error logs directory: %s\n", logDir)
}
