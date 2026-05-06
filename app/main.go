package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
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
	if _, exists := commands[cmd]; !exists {
		fmt.Printf("%s: invalid_command", cmd)
		return nil
	}

	fmt.Printf("%s is a shell buitin", cmd)
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
		if !exists {
			fmt.Printf("%s: command not found\n", cmd)
			continue
		}

		err := fn(args)
		if err == errExit {
			break outer
		}

		if err != nil {
			fmt.Printf("%v", err)
		}
	}
}
