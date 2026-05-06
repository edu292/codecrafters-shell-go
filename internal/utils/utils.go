package utils

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

var ErrNotFoundPath = errors.New("file not found in $PATH")

type parserState int

const (
	normal parserState = iota
	inDoubleQuote
	inSingleQuote
)

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

	var origin string
	relPath, found := strings.CutPrefix(relPath, "~")
	if found {
		origin, _ = os.UserHomeDir()
	} else {
		origin, _ = os.Getwd()
	}

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

func ParseInput(input []byte) (string, []string) {
	before, after, ok := bytes.Cut(input, []byte(" "))
	if !ok {
		return string(input), []string{""}
	}

	cmd := string(before)

	var args []string
	var buf []byte
	parserState := normal
	for _, b := range after {
		switch {
		case parserState == inDoubleQuote && b == '"':
			parserState = normal
		case parserState == inSingleQuote && b == '\'':
			parserState = normal
		case parserState == normal:
			switch b {
			case ' ':
				if len(buf) > 0 {
					args = append(args, string(buf))
				}
				buf = buf[:0]
			case '\'':
				parserState = inSingleQuote
			case '"':
				parserState = inDoubleQuote
			default:
				buf = append(buf, b)
			}
		default:
			buf = append(buf, b)
		}
	}

	if len(buf) > 0 {
		args = append(args, string(buf))
	}

	return cmd, args
}
