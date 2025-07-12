package user_input

import (
    "database/sql"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "strings"

    "dalennod/internal/backup"
    "dalennod/internal/constants"
    "dalennod/internal/db"
    "dalennod/internal/server"
    "dalennod/internal/setup"
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
    case setup.FlagVals.FixDB:
        applyDBUpdates(database)
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
    if err != nil {
        fmt.Println("error getting current profile directory. ERROR:", err)
        return
    }

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
    fmt.Printf("Error logs directory: %s\n", constants.LOGS_PATH)
}

func applyDBUpdates(database *sql.DB) {
    rows, err := database.Query("SELECT id, byteThumbURL from bookmarks WHERE byteThumbURL NOT NULL;")
    if err != nil {
        fmt.Println("error querying database. ERROR:", err)
        return
    }

    var id int
    var thumb []byte

    for rows.Next() {
        rows.Scan(&id, &thumb)

        writeFilePath := filepath.Join(constants.THUMBNAILS_PATH, strconv.Itoa(id))
        err := os.WriteFile(writeFilePath, thumb, 0644)
        if err != nil {
            fmt.Printf("failed to create local thumbnail for ID %d. ERROR: %v\n", id, err)
            continue
        }
    }

    if _, err := database.Exec("ALTER TABLE bookmarks DROP COLUMN byteThumbURL;"); err != nil {
        fmt.Println("failed to drop column byteThumbURL. ERROR:", err)
        fmt.Println("needs manual intervention")
        return
    }
}
