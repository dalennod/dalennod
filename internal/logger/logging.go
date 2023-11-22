package logger

import (
	"fmt"
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
	logPath            = "./logs"
)

func Enable() {
	err := os.MkdirAll(logPath, 0755)
	if err != nil {
		log.Fatalf("Error creating logs directory: %v\n", err)
	}

	logFile, err := os.Create(fmt.Sprintf("%s/%s.log", logPath, logFileName))
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
