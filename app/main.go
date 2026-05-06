package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"codecrafters-shell-go/internal/utils"
)

var errExit = errors.New("exit")

type command func(args []string) error

func exit(args []string) error {
	return errExit
}

func echo(args []string) error {
	fmt.Println(strings.Join(args, " "))

	return nil
}

func typeHandler(args []string) error {
	cmd := args[0]
	if _, exists := commands[cmd]; exists {
		fmt.Printf("%s is a shell builtin\n", cmd)
		return nil
	}

	if fullpath, err := utils.LookPath(cmd); err == nil {
		fmt.Printf("%s is %s\n", cmd, fullpath)
		return nil
	}

	fmt.Printf("%s: not found\n", cmd)
	return nil
}

var commands = make(map[string]command)

func init() {
	commands["exit"] = exit
	commands["echo"] = echo
	commands["type"] = typeHandler
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
outer:
	for {
		fmt.Print("$ ")
		if !scanner.Scan() {
			break
		}

		fields := strings.Fields(scanner.Text())
		cmd, args := fields[0], fields[1:]

		fn, exists := commands[cmd]
		if exists {
			err := fn(args)
			if err == errExit {
				break outer
			}

			if err != nil {
				fmt.Printf("%v", err)
			}

			continue
		}

		fullpath, err := utils.LookPath(cmd)
		if err != nil {
			fmt.Printf("%s: command not found\n", cmd)
			continue
		}

		ex := exec.Command(fullpath, fields...)
		ex.Stdout = os.Stdout
		ex.Stdin = os.Stdin
		ex.Stderr = os.Stderr

		err = ex.Run()
		if err != nil {
			fmt.Print(err)
		}
	}
}
