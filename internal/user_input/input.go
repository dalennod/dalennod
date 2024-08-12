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
	"encoding/json"
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

	// should display config dir path and log dir path
	// on first-run. todo in future
	// cfgDir, logDir, err := enableLogs()
	// if err != nil {
	// 	logger.Warn.Println("Could not enable error logging")
	// }

	database = data

	flagVals = setup.ParseFlags()

	switch true {
	case flagVals.ViewAll:
		db.ViewAll(database, false)
	case flagVals.StartServer:
		server.Start(database)
	case flagVals.AddEntry:
		addInput(setup.Bookmark{}, false)
	case flagVals.Backup && flagVals.JSONOut:
		backup.JSONOut(database)
	case flagVals.Where:
		whereConfigLog()
	}

	if flagVals.RemoveID != "" {
		removeInput(flagVals.RemoveID)
	} else if flagVals.UpdateID != "" {
		updateInput(flagVals.UpdateID)
	} else if flagVals.ViewID != "" {
		viewInput(flagVals.ViewID)
	} else if flagVals.Import != "" && flagVals.Firefox {
		importFirefoxInput(flagVals.Import)
	} else if flagVals.Import != "" && flagVals.Dalennod {
		importDalennodInput(flagVals.Import)
	}
}

func addInput(bmStruct setup.Bookmark, callToUpdate bool) {
	var archiveUrl string
	bmStruct, archiveUrl = getBmInfo(bmStruct)

	if !callToUpdate {
		addBm(bmStruct, archiveUrl)
	} else {
		updateBm(bmStruct, archiveUrl)
	}
}

func getBmInfo(bmStruct setup.Bookmark) (setup.Bookmark, string) {
	var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
	var err error

	fmt.Print("URL to save: ")
	scanner.Scan()
	bmStruct.URL = scanner.Text()

	bmStruct.ThumbURL, bmStruct.ByteThumbURL, err = thumb_url.GetPageThumb(bmStruct.URL)
	if err != nil {
		bmStruct.ThumbURL = bmStruct.URL
	}

	fmt.Print("Title for the bookmark: ")
	scanner.Scan()
	bmStruct.Title = scanner.Text()

	fmt.Print("Notes/log reason for bookmark: ")
	scanner.Scan()
	bmStruct.Note = scanner.Text()

	fmt.Print("Keywords for searching later: ")
	scanner.Scan()
	bmStruct.Keywords = scanner.Text()

	fmt.Print("Group to store the bookmark into: ")
	scanner.Scan()
	bmStruct.BmGroup = scanner.Text()

	fmt.Print("Archive URL? (y/N): ")
	scanner.Scan()
	var archiveUrl string = scanner.Text()

	return bmStruct, archiveUrl
}

func updateBm(bmStruct setup.Bookmark, archiveUrl string) {
	switch archiveUrl {
	case "y", "Y":
		bmStruct.Archived, bmStruct.SnapshotURL = archive.SendSnapshot(bmStruct.URL)
		if bmStruct.Archived {
			db.Update(database, bmStruct, false)
		} else {
			logger.Warn.Println("Snapshot failed")
			db.Update(database, bmStruct, false)
		}
	case "n", "N", "":
		db.Update(database, bmStruct, false)
	default:
		db.Update(database, bmStruct, false)
		logger.Warn.Println("Invalid input for archive request. URL has not been archived")
	}
}

func addBm(bmStruct setup.Bookmark, archiveUrl string) {
	switch archiveUrl {
	case "y", "Y":
		bmStruct.Archived, bmStruct.SnapshotURL = archive.SendSnapshot(bmStruct.URL)
		if bmStruct.Archived {
			db.Add(database, bmStruct)
		} else {
			logger.Warn.Println("Snapshot failed")
			db.Add(database, bmStruct)
		}
	case "n", "N", "":
		db.Add(database, bmStruct)
	default:
		db.Add(database, bmStruct)
		logger.Warn.Println("Invalid input for archive request. URL has not been archived")
	}
}

func updateInput(updateID string) {
	var (
		confirm string
		bmAtID  setup.Bookmark
		scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
	)

	idToINT, err := strconv.Atoi(updateID)
	if err != nil {
		logger.Error.Println("Invalid input")
	}

	bmAtID, err = db.ViewSingleRow(database, idToINT, false)
	if err != nil {
		fmt.Println(err)
		logger.Error.Println(err)
		return
	}
	fmt.Println(bmAtID)

	fmt.Print("Update this entry? (Y/n): ")
	scanner.Scan()
	confirm = scanner.Text()

	switch confirm {
	case "y", "Y", "":
		fmt.Println("Leave empty to retain old information")
		addInput(bmAtID, true)
	case "n", "N":
		return
	default:
		logger.Info.Println("Invalid input. Received:", confirm)
		fmt.Println("Invalid input. Exiting")
		return
	}
}

func removeInput(removeID string) {
	var (
		confirm string
		scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
	)

	idToINT, err := strconv.Atoi(removeID)
	if err != nil {
		logger.Error.Println("Invalid input")
		return
	}

	_, err = db.ViewSingleRow(database, idToINT, false)
	if err != nil {
		fmt.Println(err)
		logger.Error.Println(err)
		return
	}

	fmt.Print("Remove this entry? (Y/n): ")
	scanner.Scan()
	confirm = scanner.Text()

	switch confirm {
	case "y", "Y", "":
		db.Remove(database, idToINT)
	case "n", "N":
		return
	default:
		logger.Info.Println("Invalid input received:", confirm)
		fmt.Println("Invalid input. Exiting")
		return
	}
}

func viewInput(viewID string) {
	idToINT, err := strconv.Atoi(viewID)
	if err != nil {
		logger.Error.Println("Invalid input")
		return
	}
	_, err = db.ViewSingleRow(database, idToINT, false)
	if err != nil {
		fmt.Println(err)
		logger.Error.Println(err)
		return
	}
}

func importDalennodInput(file string) {
	fContent, err := os.ReadFile(file)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	var importedBookmarks []setup.Bookmark
	err = json.Unmarshal(fContent, &importedBookmarks)
	if err != nil {
		logger.Error.Fatalln(err)
	}
	var importedBookmarksCount = len(importedBookmarks)
	for i, importedData := range importedBookmarks {
		db.Add(database, importedData)
		fmt.Printf("\rAdded %d / %d", i+1, importedBookmarksCount)
	}
	fmt.Println()
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
		db.Add(database, parsedBookmark)
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

func whereConfigLog() {
	cfgDir, err := setup.ConfigDir()
	if err != nil {
		fmt.Println(err)
	}

	logDir, err := setup.CacheDir()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Database and config directory: %s\n", cfgDir)
	fmt.Printf("Error logs directory: %s\n", logDir)
}

func enableLogs() (string, string, error) {
	logger.Enable()

	cfgDir, err := setup.ConfigDir()
	if err != nil {
		return "", "", err
	}

	logDir, err := setup.CacheDir()
	if err != nil {
		return "", "", err
	}

	logger.Info.Printf("Database and config directory: %s\n", cfgDir)
	logger.Info.Printf("Error logs directory: %s\n", logDir)

	return cfgDir, logDir, nil
}
