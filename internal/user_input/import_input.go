package user_input

import (
	"dalennod/internal/bookmark_import"
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"dalennod/internal/setup"
	"encoding/json"
	"fmt"
	"os"
)

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

func importFirefoxInput(file string) {
	readFile, err := os.Open(file)
	if err != nil {
		logger.Error.Printf("couldn't open file. ERROR: %v", err)
	}

	firefoxBookmarks := &bookmark_import.Item{}
	if err := json.NewDecoder(readFile).Decode(firefoxBookmarks); err != nil {
		logger.Error.Println(err)
		fmt.Println(err)
		return
	}

	parsedBookmarks := bookmark_import.ParseFirefox(firefoxBookmarks, "")
	parsedBookmarksLength := len(parsedBookmarks)

	for i, parsedBookmark := range parsedBookmarks {
		db.Add(database, parsedBookmark)
		fmt.Printf("\rAdded %d / %d", i+1, parsedBookmarksLength)
	}
	fmt.Println("")
}

func importChromiumInput(file string) {
	readFile, err := os.Open(file)
	if err != nil {
		logger.Error.Printf("couldn't open file. ERROR: %v", err)
	}

	var chromiumBookmarks bookmark_import.ChromiumItem
	if err := json.NewDecoder(readFile).Decode(&chromiumBookmarks); err != nil {
		logger.Error.Println(err)
		fmt.Println(err)
		return
	}

	parsedBookmarks := bookmark_import.ParseChromium(chromiumBookmarks)
	parsedBookmarksLength := len(parsedBookmarks)

	for i, parsedBookmark := range parsedBookmarks {
		db.Add(database, parsedBookmark)
		fmt.Printf("\rAdded %d / %d", i+1, parsedBookmarksLength)
	}
	fmt.Println("")
}
