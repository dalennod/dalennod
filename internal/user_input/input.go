package user_input

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"dalennod/internal/archive"
	"dalennod/internal/backup"
	"dalennod/internal/constants"
	"dalennod/internal/db"
	"dalennod/internal/server"
	"dalennod/internal/setup"
	"dalennod/internal/thumb_url"
)

var database *sql.DB

func UserInput(bookmark_database *sql.DB) {
	database = bookmark_database

	switch {
	case setup.FlagVals.StartServer:
		server.Start(database)
	case setup.FlagVals.Where:
		whereConfigLog()
	case setup.FlagVals.ViewAll:
		viewAllInput(db.ViewAll(database))
	case setup.FlagVals.ViewID != "":
		viewInput(setup.FlagVals.ViewID)
	case setup.FlagVals.RemoveID != "":
		removeInput(setup.FlagVals.RemoveID)
	case setup.FlagVals.AddEntry:
		addInput(setup.Bookmark{}, false)
	case setup.FlagVals.UpdateID != "":
		updateInput(setup.FlagVals.UpdateID)
	case setup.FlagVals.Profile && setup.FlagVals.Switch == "":
		showProfiles()
	case setup.FlagVals.Profile && setup.FlagVals.Switch != "":
		switchProfile(setup.FlagVals.Switch)
	case setup.FlagVals.Backup && setup.FlagVals.JSONOut:
		backup.JSONOut(database)
	case setup.FlagVals.RedoCompletion:
		setup.SetCompletion()
	case setup.FlagVals.Import && setup.FlagVals.Firefox != "":
		importFirefoxInput(setup.FlagVals.Firefox)
	case setup.FlagVals.Import && setup.FlagVals.Chromium != "":
		importChromiumInput(setup.FlagVals.Chromium)
	case setup.FlagVals.Import && setup.FlagVals.Dalennod != "":
		importDalennodInput(setup.FlagVals.Dalennod)
	// case setup.FlagVals.FixDB:
	// 	applyDBUpdates(database)
	}
}

func showProfiles() {
	osConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("error getting config directory. ERROR:", err)
		return
	}

	openConfigDir, err := os.Open(osConfigDir)
	if err != nil {
		fmt.Println("error opening config dir. ERROR:", err)
		return
	}
	defer openConfigDir.Close()

	dirEntries, err := openConfigDir.Readdir(-1)
	if err != nil {
		fmt.Println("error reading directory entries. ERROR:", err)
		return
	}

	for _, entry := range dirEntries {
		if entry.IsDir() && strings.Contains(entry.Name(), constants.NAME) {
			profileName := strings.SplitN(entry.Name(), ".", 2)
			fmt.Println(printProfileNames(profileName))
		}
	}
}

func printProfileNames(profileName []string) string {
	if len(profileName) > 1 {
		return "  " + profileName[1]
	} else {
		return "* current"
	}
}

func switchProfile(profileName string) {
	osConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("error getting config directory. ERROR:", err)
		return
	}

	switchProfilePath := filepath.Join(osConfigDir, constants.NAME+"."+profileName)
	if _, err := os.Stat(switchProfilePath); os.IsNotExist(err) {
		fmt.Printf("profile \"%s\" does not exist. ERROR: %v\n", profileName, err)
		return
	}

	currentProfileDir := constants.CONFIG_PATH

	var userInput string
	fmt.Printf("Final rename will be \"%s.{your_input}\"\nRename current profile to: ", constants.NAME)
	_, err = fmt.Scanln(&userInput)
	if err != nil {
		fmt.Println("error reading input. ERROR:", err)
		return
	}

	currentRename := filepath.Join(osConfigDir, constants.NAME+"."+userInput)
	if err := os.Rename(currentProfileDir, currentRename); err != nil {
		fmt.Println("error renaming current profile. ERROR:", err)
		return
	}

	if err := os.Rename(switchProfilePath, currentProfileDir); err != nil {
		fmt.Println("error switching profile. ERROR:", err)
		return
	}

	fmt.Printf("OLD: %s -> %s\nNEW: %s -> %s\nProfile switched.\n", currentProfileDir, currentRename, switchProfilePath, currentProfileDir)
}

func whereConfigLog() {
	fmt.Printf("Database and config directory: %s\n", constants.CONFIG_PATH)
}

// func applyDBUpdates(database *sql.DB) {
// 	rows, err := database.Query("SELECT id, byteThumbURL from bookmarks WHERE byteThumbURL NOT NULL;")
// 	if err != nil {
// 		fmt.Println("error querying database. ERROR:", err)
// 		return
// 	}

// 	var id int
// 	var thumb []byte

// 	for rows.Next() {
// 		rows.Scan(&id, &thumb)

// 		writeFilePath := filepath.Join(constants.THUMBNAILS_PATH, strconv.Itoa(id))
// 		err := os.WriteFile(writeFilePath, thumb, 0644)
// 		if err != nil {
// 			fmt.Printf("failed to create local thumbnail for ID %d. ERROR: %v\n", id, err)
// 			continue
// 		}
// 	}

// 	if _, err := database.Exec("ALTER TABLE bookmarks DROP COLUMN byteThumbURL;"); err != nil {
// 		fmt.Println("failed to drop column byteThumbURL. ERROR:", err)
// 		fmt.Println("needs manual intervention")
// 		return
// 	}

// 	if _, err := database.Exec("VACUUM;"); err != nil {
// 		fmt.Println("error running 'VACUUM;' to rebuild database file & reduce database size")
// 		fmt.Println("does not need manual intervention, but recommended")
// 		return
// 	}
// }

func addInput(bkmStruct setup.Bookmark, callToUpdate bool) {
	var archiveUrl string
	bkmStruct, archiveUrl = getBKMInfo(bkmStruct)

	if !callToUpdate {
		addBKM(bkmStruct, archiveUrl)
	} else {
		updateBKM(bkmStruct, archiveUrl)
	}
}

func getBKMInfo(bkmStruct setup.Bookmark) (setup.Bookmark, string) {
	var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
	var err error

	fmt.Print("URL to save: ")
	scanner.Scan()
	bkmStruct.URL = scanner.Text()

	bkmStruct.ThumbURL, err = thumb_url.GetPageThumb(bkmStruct.URL)
	if err != nil {
		bkmStruct.ThumbURL = bkmStruct.URL
	}

	fmt.Print("Title for the bookmark: ")
	scanner.Scan()
	bkmStruct.Title = scanner.Text()

	fmt.Print("Notes/log reason for bookmark: ")
	scanner.Scan()
	bkmStruct.Note = scanner.Text()

	fmt.Print("Keywords for searching later: ")
	scanner.Scan()
	bkmStruct.Keywords = scanner.Text()

	fmt.Print("Category to store the bookmark into: ")
	scanner.Scan()
	bkmStruct.Category = scanner.Text()

	fmt.Print("Archive URL? (y/N): ")
	scanner.Scan()
	var archiveUrl string = scanner.Text()

	return bkmStruct, archiveUrl
}

func updateBKM(bkmStruct setup.Bookmark, archiveUrl string) {
	switch archiveUrl {
	case "y", "Y":
		bkmStruct.Archived, bkmStruct.SnapshotURL = archive.SendSnapshot(bkmStruct.URL)
		if bkmStruct.Archived {
			db.Update(database, bkmStruct, false)
		} else {
			log.Println("WARN: snapshot failed")
			db.Update(database, bkmStruct, false)
		}
	case "n", "N", "":
		db.Update(database, bkmStruct, false)
	default:
		db.Update(database, bkmStruct, false)
		log.Println("WARN: invalid input for archive request. URL has not been archived. got input:", archiveUrl)
	}
}

func addBKM(bkmStruct setup.Bookmark, archiveUrl string) {
	switch archiveUrl {
	case "y", "Y":
		bkmStruct.Archived, bkmStruct.SnapshotURL = archive.SendSnapshot(bkmStruct.URL)
		if bkmStruct.Archived {
			db.Add(database, bkmStruct)
		} else {
			log.Println("WARN: snapshot failed")
			db.Add(database, bkmStruct)
		}
	case "n", "N", "":
		db.Add(database, bkmStruct)
	default:
		db.Add(database, bkmStruct)
		log.Println("WARN: invalid input for archive request. URL has not been archived. got input:", archiveUrl)
	}
}

func updateInput(updateID string) {
	var (
		confirm string
		scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
	)

	idToINT, err := strconv.ParseInt(updateID, 10, 64)
	if err != nil {
		log.Println("WARN: invalid bookmark ID:", err)
		return
	}

	bkmAtID, err := db.ViewSingleRow(database, idToINT)
	if err != nil {
		log.Printf("WARN: could not get record for bookmark id: %s: %v\n", updateID, err)
		return
	}
	db.PrintRow(bkmAtID)

	fmt.Print("Update this entry? (Y/n): ")
	scanner.Scan()
	confirm = scanner.Text()

	switch confirm {
	case "y", "Y", "":
		fmt.Println("Leave empty to retain old information")
		addInput(bkmAtID, true)
	case "n", "N":
		return
	default:
		log.Println("WARN: invalid input:", confirm)
		return
	}
}

func removeInput(removeID string) {
	var (
		confirm string
		scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
	)

	idToINT, err := strconv.ParseInt(removeID, 10, 64)
	if err != nil {
		log.Println("WARN: Invalid input:", err)
		return
	}

	bkmAtID, err := db.ViewSingleRow(database, idToINT)
	if err != nil {
		log.Printf("WARN: could not get record for bookmark id: %s: %v\n", removeID, err)
		return
	}
	db.PrintRow(bkmAtID)

	fmt.Print("Remove this entry? (Y/n): ")
	scanner.Scan()
	confirm = scanner.Text()

	switch confirm {
	case "y", "Y", "":
		db.Remove(database, idToINT)
	case "n", "N":
		return
	default:
		log.Println("WARN: invalid input:", confirm)
		return
	}
}

func viewInput(viewID string) {
	idToINT, err := strconv.ParseInt(viewID, 10, 64)
	if err != nil {
		log.Println("Invalid input")
		return
	}

	bkmAtID, err := db.ViewSingleRow(database, idToINT)
	if err != nil {
		log.Printf("WARN: could not get record for bookmark id: %s: %v\n", viewID, err)
		return
	}

	db.PrintRow(bkmAtID)
}

func viewAllInput(bookmarks []setup.Bookmark) {
	if len(bookmarks) == 0 {
		log.Println("WARN: database empty when trying to view all")
		return
	}

	for _, bookmark := range bookmarks {
		db.PrintRow(bookmark)
	}
}
