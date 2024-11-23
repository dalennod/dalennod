package logger

import (
	"dalennod/internal/setup"
	"io"
	"log"
	"os"
	"time"
)

var (
	Warn  *log.Logger
	Info  *log.Logger
	Error *log.Logger

	logFileName string = time.Now().Format("2006-01-02_15.04.05")
)

func Enable() {
	logPath, err := setup.CacheDir()
	if err != nil {
		log.Fatalln(err)
	}

	logFile, err := os.Create(logPath + "/" + logFileName + ".log")
	if err != nil {
		log.Fatalln("Error creating log file:", err)
	}

	Warn = log.New(logFile, "WARN: ", log.Lmicroseconds|log.Lshortfile)
	Info = log.New(logFile, "INFO: ", log.Lmicroseconds|log.Lshortfile)
	Error = log.New(logFile, "ERROR: ", log.Lmicroseconds|log.Lshortfile)
}

func Disable() {
	Warn.SetOutput(io.Discard)
	Info.SetOutput(io.Discard)
	Error.SetOutput(io.Discard)
}
