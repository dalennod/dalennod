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

var (
	database *sql.DB
	flagVals setup.FlagValues
)

func UserInput(data *sql.DB) {
	enableLogs()
	database = data

	flagVals = setup.ParseFlags()
	switch {
	case flagVals.StartServer:
		server.Start(database)
	case flagVals.Where:
		whereConfigLog()
	case flagVals.ViewAll:
		db.ViewAll(database, false)
	case flagVals.ViewID != "":
		viewInput(flagVals.ViewID)
	case flagVals.RemoveID != "":
		removeInput(flagVals.RemoveID)
	case flagVals.AddEntry:
		addInput(setup.Bookmark{}, false)
	case flagVals.UpdateID != "":
		updateInput(flagVals.UpdateID)
	case flagVals.Profile && flagVals.Switch == "":
		showProfiles()
	case flagVals.Profile && flagVals.Switch != "":
		switchProfile(flagVals.Switch)
	case flagVals.Backup && flagVals.JSONOut:
		backup.JSONOut(database)
	case flagVals.Import && flagVals.Firefox != "":
		importFirefoxInput(flagVals.Firefox)
	case flagVals.Import && flagVals.Chromium != "":
		importChromiumInput(flagVals.Chromium)
	case flagVals.Import && flagVals.Dalennod != "":
		importDalennodInput(flagVals.Dalennod)
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
