package utils

import (
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

var ErrNotFoundPath = errors.New("file not found in $PATH")

func LookPath(cmd string) (string, error) {
	path := os.Getenv("PATH")
	pathDirs := filepath.SplitList(path)

	for _, dir := range pathDirs {
		fullPath := filepath.Join(dir, cmd)

		if err := unix.Access(fullPath, unix.X_OK); err != nil {
			continue
		}

		return fullPath, nil
	}
	return "", ErrNotFoundPath
}
