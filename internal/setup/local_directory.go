package setup

import (
	"dalennod/internal/logger"

	"os"
	"runtime"
)

const (
	NAME = "dalennod/"

	UNIX_PATH = "/.local/share/" + NAME
	WIN_PATH  = "/AppData/Roaming/" + NAME
)

func GetOS() string {
	logger.Enable()

	var os string = runtime.GOOS
	var path string

	homeDir, err := getHomeDir()
	if err != nil {
		logger.Error.Fatalln("Could not locate home directory", err)
	}

	switch os {
	case "linux", "darwin":
		path = unixSetup(homeDir)
	case "windows":
		path = winSetup(homeDir)
	default:
		logger.Error.Fatalln("Failed to recognize OS:", os)
	}

	return path
}

func getHomeDir() (string, error) {
	return os.UserHomeDir()
}

func winSetup(homeDir string) string {
	var path string = homeDir + WIN_PATH
	createDir(path)
	return path
}

func unixSetup(homeDir string) string {
	var path string = homeDir + UNIX_PATH
	createDir(path)
	return path
}

func createDir(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		logger.Error.Fatalf("Error creating database directory: %v\n", err)
	}
	logger.Info.Printf("Database directory created at %s\n", path)
}
