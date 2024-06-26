package setup

import (
	"log"

	"os"
	"runtime"
)

const (
	NAME string = "/dalennod/"
	LOGS string = "logs/"
	DB   string = "db/"
)

func GetOS() string {
	cfgDir, err := ConfigDir()
	if err != nil {
		log.Fatalln(err)
	}

	cacheDir, err := CacheDir()
	if err != nil {
		log.Fatalln(err)
	}

	dbDir, err := DatabaseDir()
	if err != nil {
		log.Fatalln(err)
	}

	var goos string = runtime.GOOS

	switch goos {
	case "linux", "darwin":
		createDir(cfgDir, dbDir, cacheDir)
		// TODO: fish (& bash) completion did not work. generated completion using .fish file did not work properly. also, was not able to run individual complete commands using exec() because 'complete' was not recognized. try again later
		// defer SetCompletion()
	case "windows":
		createDir(cfgDir, dbDir, cacheDir)
	default:
		log.Fatalln("unrecognized OS:", err)
	}

	CfgSetup(cfgDir)

	return dbDir
}

func ConfigDir() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return cfgDir + NAME, nil
}

func CacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return cacheDir + NAME + LOGS, nil
}

func DatabaseDir() (string, error) {
	dbDir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return dbDir + DB, nil
}

func createDir(args ...string) {
	for _, path := range args {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
