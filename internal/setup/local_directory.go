package setup

import (
    "dalennod/internal/constants"
    "dalennod/internal/logger"
    "log"
    "os"
    "runtime"
    "path/filepath"
)

func setupDirectories() string {
    if _, err := configDir(); err != nil {
        log.Fatalln("error getting config directory. ERROR:", err)
    }

    if _, err := databaseDir(); err != nil {
        log.Fatalln("error getting database directory. ERROR:", err)
    }

    if _, err := cacheDir(); err != nil {
        log.Fatalln("error getting logs directory. ERROR:", err)
    }

    goos := runtime.GOOS
    switch goos {
    case "linux", "darwin":
        createDir(constants.CONFIG_PATH, constants.DB_PATH, constants.LOGS_PATH)
        defer setCompletion()
    case "windows":
        createDir(constants.CONFIG_PATH, constants.DB_PATH, constants.LOGS_PATH)
    default:
        log.Fatalln("ERROR: unrecognized OS:", goos)
    }

    configSetup(constants.CONFIG_PATH)

    return constants.DB_PATH
}

func InitLocalDirs() string {
    if _, err := configDir(); err != nil {
        log.Fatalln("error getting config directory. ERROR:", err)
    }

    if _, err := databaseDir(); err != nil {
        log.Fatalln("error getting database directory. ERROR:", err)
    }
    if _, err := os.Stat(constants.DB_PATH); os.IsNotExist(err) {
        constants.DB_PATH = setupDirectories()
    } else {
        readConfig, err := ReadCfg()
        if err != nil {
            log.Fatalln("error reading config. ERROR:", err)
        }
        if readConfig.FirstRun {
            writeCfg(false)
        }
    }

    if _, err := cacheDir(); err != nil {
        log.Fatalln("error getting cache directory. ERROR:", err)
    }
    if _, err := os.Stat(constants.LOGS_PATH); os.IsNotExist(err) {
        setupDirectories()
    }

    enableLogs()

    return constants.DB_PATH
}

func enableLogs() {
    logger.Enable()

    readConfig, err := ReadCfg()
    if err != nil {
        log.Fatalln("error reading config. ERROR:", err)
    }

    if readConfig.FirstRun {
        logger.Info.Printf("Database and config directory: %s\n", constants.CONFIG_PATH)
        logger.Info.Printf("Error logs directory: %s\n", constants.LOGS_PATH)
    }
}

func configDir() (string, error) {
    cfgDir, err := os.UserConfigDir()
    if err != nil {
        return "", err
    }
    constants.CONFIG_PATH = filepath.Join(cfgDir, constants.NAME)
    return constants.CONFIG_PATH, nil
}

func cacheDir() (string, error) {
    cacheDir, err := os.UserCacheDir()
    if err != nil {
        return "", err
    }
    constants.LOGS_PATH = filepath.Join(cacheDir, constants.NAME, constants.LOGS_DIRNAME)
    return constants.LOGS_PATH, nil
}

func databaseDir() (string, error) {
    constants.DB_PATH = filepath.Join(constants.CONFIG_PATH, constants.DB_DIRNAME)
    return constants.DB_PATH, nil
}

func createDir(args ...string) {
    for _, path := range args {
        err := os.MkdirAll(path, 0755)
        if err != nil {
            log.Fatalln("error creating directories. ERROR:", err)
        }
    }
}
