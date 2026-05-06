package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

var ErrNotFoundPath = errors.New("file not found in $PATH")

func GetFromPath(cmd string) (string, error) {
	path := os.Getenv("PATH")
	pathDirs := filepath.SplitList(path)

	for _, dir := range pathDirs {
		absolutePath := filepath.Join(dir, cmd)

		if err := unix.Access(absolutePath, unix.X_OK); err != nil {
			continue
		}

		return absolutePath, nil
	}
	return "", ErrNotFoundPath
}

func ExpandPath(relPath string) string {
	origin, _ := os.Getwd()

	relPath, _ = strings.CutPrefix(relPath, "./")
	for strings.HasPrefix(relPath, "..") {
		origin = filepath.Base(origin)
	}

	return filepath.Join(origin, relPath)
}
