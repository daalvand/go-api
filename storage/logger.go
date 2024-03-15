package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

func InitLogger(params ...string) {
	var logDir, logFileName string

	if len(params) > 0 {
		logDir = params[0]
	}
	if len(params) > 1 {
		logFileName = params[1]
	}

	if logDir == "" {
		logDir = getDefaultLogDir()
	}

	if logFileName == "" {
		logFileName = generateDefaultLogFileName()
	}

	absLogDir, err := filepath.Abs(logDir)
	if err != nil {
		logrus.Fatalf("Error getting absolute path for log directory: %s", err)
	}

	if err := os.MkdirAll(absLogDir, 0755); err != nil {
		logrus.Fatalf("Error creating log directory: %s", err)
	}

	logFilePath := filepath.Join(absLogDir, logFileName)

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalf("Error opening log file: %s", err)
	}

	logrus.SetOutput(logFile)

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
}

func getDefaultLogDir() string {
	_, currentFile, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(currentFile), "logs")
}

func generateDefaultLogFileName() string {
	return fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
}
