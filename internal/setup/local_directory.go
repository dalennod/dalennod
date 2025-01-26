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

	dbDir, err := DatabaseDir()
	if err != nil {
		log.Fatalln(err)
	}

	if !CfgSetup(cfgDir) {
		return dbDir
	}

	cacheDir, err := CacheDir()
	if err != nil {
		log.Fatalln(err)
	}

	var goos string = runtime.GOOS

	switch goos {
	case "linux", "darwin":
		createDir(cfgDir, dbDir, cacheDir)
		defer setCompletion()
	case "windows":
		createDir(cfgDir, dbDir, cacheDir)
	default:
		log.Fatalln("unrecognized OS:", err)
	}

	return dbDir
}

func ConfigDir() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return cfgDir + NAME, nil // start to use filepath.Join
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
