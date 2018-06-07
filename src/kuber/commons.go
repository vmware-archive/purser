package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

func executeCommand(command string) []byte {
	slice := strings.Fields(command)
	cmd := exec.Command(slice[0], slice[1:]...)
	cmd.Env = os.Environ()
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return out.Bytes()
}
