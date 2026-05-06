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
	if strings.HasPrefix(relPath, "/") {
		return relPath
	}

	origin, _ := os.Getwd()

	var found bool
	for {
		relPath, found = strings.CutPrefix(relPath, "..")
		if !found {
			break
		}

		origin = filepath.Dir(origin)
	}
	relPath, _ = strings.CutPrefix(relPath, ".")

	return filepath.Join(origin, relPath)
}
