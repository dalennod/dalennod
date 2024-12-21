package user_input

import (
	"dalennod/internal/backup"
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"dalennod/internal/server"
	"dalennod/internal/setup"
	"database/sql"
	"fmt"
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
	} else if flagVals.Import != "" && flagVals.Chromium {
		importChromiumInput(flagVals.Import)
	} else if flagVals.Import != "" && flagVals.Dalennod {
		importDalennodInput(flagVals.Import)
	}
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
