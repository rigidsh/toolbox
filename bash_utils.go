package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func which(cmd string) (string, bool) {

	for _, path := range strings.Split(os.Getenv("PATH"), ":") {
		cmdPath := filepath.Join(path, cmd)

		_, err := os.Stat(cmdPath)
		if !os.IsNotExist(err) {
			return cmdPath, true
		}
	}

	return "", false
}

func cp(from, to string) error {
	copyCommand := exec.Command("cp", "-r", from, to)

	return copyCommand.Run()
}
