package utils

import (
	"os"
	"os/user"

	log "github.com/Sirupsen/logrus"
)

// OpenFile handles opening file in Read/Write mode, creating and appending to it as needed.
func OpenFile(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		log.Errorf("failed to open file %s, %v", filename, err)
	}
	return f
}

// GetUsrHomeDir returns the current user's Home Directory
func GetUsrHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Errorf("failed to fetch current user %v", err)
	}
	return usr.HomeDir
}
