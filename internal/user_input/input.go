package user_input

import (
    "dalennod/internal/backup"
    "dalennod/internal/db"
    "dalennod/internal/logger"
    "dalennod/internal/server"
    "dalennod/internal/setup"
    "database/sql"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

var database *sql.DB

func UserInput(bookmark_database *sql.DB) {
    enableLogs()
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
    case setup.FlagVals.Import && setup.FlagVals.Firefox != "":
        importFirefoxInput(setup.FlagVals.Firefox)
    case setup.FlagVals.Import && setup.FlagVals.Chromium != "":
        importChromiumInput(setup.FlagVals.Chromium)
    case setup.FlagVals.Import && setup.FlagVals.Dalennod != "":
        importDalennodInput(setup.FlagVals.Dalennod)
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
    }

    for _, entry := range dirEntries {
        if entry.IsDir() && strings.Contains(entry.Name(), "dalennod") {
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

    switchProfilePath := filepath.Join(osConfigDir, "dalennod."+profileName)
    if _, err := os.Stat(switchProfilePath); os.IsNotExist(err) {
        fmt.Printf("profile \"%s\" does not exist. ERROR: %v\n", profileName, err)
        return
    }

    currentProfileDir, err := setup.ConfigDir()
    if err != nil {
        fmt.Println("error getting current profile directory. ERROR:", err)
        return
    }

    var userInput string
    fmt.Print("Final rename will be \"dalennod.{your_input}\"\nRename current profile to: ")
    _, err = fmt.Scanln(&userInput)
    if err != nil {
        fmt.Println("error reading input. ERROR:", err)
        return
    }

    currentRename := filepath.Join(osConfigDir, "dalennod."+userInput)
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

func enableLogs() {
    logger.Enable()

    cfgDir, err := setup.ConfigDir()
    if err != nil {
        logger.Warn.Println("error getting config directory. ERROR:", err)
    }

    logDir, err := setup.CacheDir()
    if err != nil {
        logger.Info.Println("error getting cache directory. ERROR:", err)
    }

    logger.Info.Printf("Database and config directory: %s\n", cfgDir)
    logger.Info.Printf("Error logs directory: %s\n", logDir)
}
