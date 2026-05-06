package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

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

func argByteToString(arg []byte) string {
	return string(bytes.ReplaceAll(arg, []byte("'"), []byte("")))
}

func ParseInput(input []byte) (string, []string) {
	before, after, ok := bytes.Cut(input, []byte(" "))
	if !ok {
		return string(input), []string{}
	}

	cmd := string(before)
	byteArgs := after

	var args []string
	var startParseIdx int
	var idx int
	for idx < len(byteArgs) {
		startParseIdx = bytes.IndexFunc(byteArgs[idx:], func(r rune) bool {
			return !unicode.IsSpace(r)
		})
		if startParseIdx == -1 {
			break
		}
		startParseIdx += idx

		if bytes.HasPrefix(byteArgs[startParseIdx:], []byte("'")) {
			startParseIdx++
			idx = bytes.Index(byteArgs[startParseIdx:], []byte("' "))
		} else {
			idx = bytes.Index(byteArgs[startParseIdx:], []byte(" "))
		}

		if idx == -1 {
			idx = len(byteArgs)
		} else {
			idx += startParseIdx
		}

		if idx-startParseIdx == 1 {
			continue
		}

		args = append(args, argByteToString(byteArgs[startParseIdx:idx]))
		fmt.Println(startParseIdx)
		fmt.Println(idx)
		fmt.Println(args)
		fmt.Println()
		time.Sleep(200 * time.Millisecond)
		idx ++

	}
	fmt.Println(startParseIdx)
	fmt.Println(idx)
	fmt.Println(args)

	return cmd, args
}
