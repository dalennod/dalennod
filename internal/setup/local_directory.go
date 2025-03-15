package setup

import (
    "log"
    "os"
    "runtime"
    "path/filepath"
)

const (
    NAME string = "dalennod"
    LOGS string = "logs"
    DB   string = "db"
)

func setupDirectories() string {
    cfgDir, err := ConfigDir()
    if err != nil {
        log.Fatalln(err)
    }

    dbDir, err := DatabaseDir()
    if err != nil {
        log.Fatalln(err)
    }

    cacheDir, err := CacheDir()
    if err != nil {
        log.Fatalln(err)
    }

    goos := runtime.GOOS
    switch goos {
    case "linux", "darwin":
        createDir(cfgDir, dbDir, cacheDir)
        defer setCompletion()
    case "windows":
        createDir(cfgDir, dbDir, cacheDir)
    default:
        log.Fatalln("ERROR: unrecognized OS:", goos)
    }

    configSetup(cfgDir)

    return dbDir
}

func InitLocalDirs() string {
    databaseDir, err := DatabaseDir()
    if err != nil {
        log.Fatalln("error getting database directory. ERROR:", err)
    }
    if _, err := os.Stat(databaseDir); os.IsNotExist(err) {
        databaseDir = setupDirectories()
    } else {
        readConfig, err := ReadCfg()
        if err != nil {
            log.Fatalln("error reading config. ERROR:", err)
        }
        if readConfig.FirstRun {
            WriteCfg(false)
        }
    }

    cacheDir, err := CacheDir()
    if err != nil {
        log.Fatalln("error getting cache directory. ERROR:", err)
    }
    if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
        setupDirectories()
    }

    return databaseDir
}

func ConfigDir() (string, error) {
    cfgDir, err := os.UserConfigDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(cfgDir, NAME), nil
}

func CacheDir() (string, error) {
    cacheDir, err := os.UserCacheDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(cacheDir, NAME, LOGS), nil
}

func DatabaseDir() (string, error) {
    dbDir, err := ConfigDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(dbDir, DB), nil
}

func createDir(args ...string) {
    for _, path := range args {
        err := os.MkdirAll(path, 0755)
        if err != nil {
            log.Fatalln(err)
        }
    }
}
