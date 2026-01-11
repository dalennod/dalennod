package setup

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"dalennod/internal/constants"
)

func InitLocalDirs() string {
	runtimeOS := runtime.GOOS
	switch runtimeOS {
	case "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln("ERROR: could not get home dir:", err)
		}
		constants.DATA_PATH = filepath.Join(homeDir, ".local", "share", constants.NAME)
	case "darwin", "windows":
		cfgDataDir, err := os.UserConfigDir()
		if err != nil {
			log.Fatalln("ERROR: could not get config/data dir:", err)
		}
		constants.DATA_PATH = filepath.Join(cfgDataDir, constants.NAME)
	default:
		log.Fatalln("ERROR: unrecognized OS:", runtimeOS)
	}

	constants.DB_PATH = databaseDir()
	constants.THUMBNAILS_PATH = thumbnailDataDir()

	checkDirsExistence(constants.DATA_PATH, constants.DB_PATH, constants.THUMBNAILS_PATH)
	checkSetEnvVar()

	return constants.DB_PATH
}

func checkDirsExistence(dirPaths ...string) {
	for _, path := range dirPaths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			continue
		}
		if path == constants.DATA_PATH {
			defer SetCompletion()
		}
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalln("error creating directories. ERROR:", err)
		}
	}
}

func checkSetEnvVar() {
	constants.WEBUI_ADDR = os.Getenv("DALENNOD_ADDR")
	if constants.WEBUI_ADDR == "" {
		constants.WEBUI_ADDR = constants.WEBUI_PORT
	}
}

func databaseDir() string {
	return filepath.Join(constants.DATA_PATH, constants.DB_DIRNAME)
}

func thumbnailDataDir() string {
	return filepath.Join(constants.DATA_PATH, constants.THUMBNAILS_DIRNAME)
}
