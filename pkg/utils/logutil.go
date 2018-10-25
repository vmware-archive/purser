package utils

import (
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
)

const logFile = "purser.log"

// InitializeLogger sets and configures logger options.
func InitializeLogger() {
	logFile := OpenFile(logFile)

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetLevel(log.InfoLevel)
}
