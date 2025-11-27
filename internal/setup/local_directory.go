package setup

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"dalennod/internal/constants"
)

func setupDirectories() string {
	if _, err := configDir(); err != nil {
		log.Fatalln("ERROR: getting config directory:", err)
	}
	databaseDir()
	thumbnailDataDir()

	goos := runtime.GOOS
	switch goos {
	case "linux", "darwin":
		createDir(constants.CONFIG_PATH, constants.DB_PATH, constants.THUMBNAILS_PATH)
		defer SetCompletion()
	case "windows":
		createDir(constants.CONFIG_PATH, constants.DB_PATH, constants.THUMBNAILS_PATH)
	default:
		log.Fatalln("ERROR: unrecognized OS:", goos)
	}

	configSetup(constants.CONFIG_PATH)

	return constants.DB_PATH
}

func InitLocalDirs() string {
	if _, err := configDir(); err != nil {
		log.Fatalln("ERROR: getting config directory:", err)
	}

	databaseDir()
	if _, err := os.Stat(constants.DB_PATH); os.IsNotExist(err) {
		constants.DB_PATH = setupDirectories()
	} else {
		readConfig, err := readCfg()
		if err != nil {
			log.Fatalln("ERROR: reading config:", err)
		}
		if readConfig.FirstRun {
			writeCfg(false, readConfig.Host, readConfig.Port)
		}
	}
	thumbnailDataDir()
	if _, err := os.Stat(constants.THUMBNAILS_PATH); os.IsNotExist(err) {
		setupDirectories()
	}

	readConfig, err := readCfg()
	if err != nil {
		log.Fatalln("ERROR: reading config:", err)
	}

	constants.WEBUI_ADDR = readConfig.Host + readConfig.Port
	if len(constants.WEBUI_ADDR) == 0 {
		log.Fatalln("ERROR: Improper config. Expected Host and Port information. Got:", constants.WEBUI_ADDR)
	}

	return constants.DB_PATH
}

func configDir() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	constants.CONFIG_PATH = filepath.Join(cfgDir, constants.NAME)
	return constants.CONFIG_PATH, nil
}

func databaseDir() string {
	constants.DB_PATH = filepath.Join(constants.CONFIG_PATH, constants.DB_DIRNAME)
	return constants.DB_PATH
}

func thumbnailDataDir() string {
	constants.THUMBNAILS_PATH = filepath.Join(constants.CONFIG_PATH, constants.THUMBNAILS_DIRNAME)
	return constants.THUMBNAILS_PATH
}

func createDir(args ...string) {
	for _, path := range args {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatalln("error creating directories. ERROR:", err)
		}
	}
}
