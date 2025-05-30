package logger

import (
    "dalennod/internal/constants"
    "errors"
    "io"
    "log"
    "os"
    "path/filepath"
)

var (
    Info  *log.Logger
    Warn  *log.Logger
    Error *log.Logger
)

func Enable() {
    var logFile *os.File
    defer logFile.Close()

    fullLogPath := filepath.Join(constants.LOGS_PATH, constants.LOGS_FILENAME)

    if _, err := os.Stat(fullLogPath); errors.Is(err, os.ErrNotExist) {
        logFile = createLogFile(fullLogPath)
    } else if checkLogSize(fullLogPath) {
        logFile = createLogFile(fullLogPath)
    } else {
        logFile, err = os.OpenFile(fullLogPath, os.O_RDWR|os.O_APPEND, 0666)
        if err != nil {
            log.Fatalln("error opening log file. ERROR:", err)
        }
    }

    Info = log.New(logFile, "INFO: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
    Warn = log.New(logFile, "WARN: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
    Error = log.New(logFile, "ERROR: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

func checkLogSize(fullLogPath string) bool {
    logFile, err := os.OpenFile(fullLogPath, os.O_RDWR|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("error opening log file. ERROR:", err)
    }
    defer logFile.Close()

    fileStat, err := logFile.Stat()
    if err != nil {
        log.Fatalln("error opening log file. ERROR:", err)
    }

    if fileStat.Size() >= constants.LOG_FILE_SIZE {
        os.Remove(fullLogPath)
        return true
    }
    return false
}

func createLogFile(fullLogPath string) *os.File {
    logFile, err := os.Create(fullLogPath)
    if err != nil {
        log.Fatalln("error creating log file. ERROR:", err)
    }
    return logFile
}

func Disable() {
    Warn.SetOutput(io.Discard)
    Info.SetOutput(io.Discard)
    Error.SetOutput(io.Discard)
}
