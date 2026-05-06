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

type handler func(args []string) error

func exitHandler(args []string) error {
	return errExit
}

func echoHandler(args []string) error {
	fmt.Println(strings.Join(args, " "))

	return nil
}

func typeHandler(args []string) error {
	cmd := args[0]
	if _, exists := builtins[cmd]; exists {
		fmt.Printf("%s is a shell builtin\n", cmd)
		return nil
	}

	if absPath, err := utils.GetFromPath(cmd); err == nil {
		fmt.Printf("%s is %s\n", cmd, absPath)
		return nil
	}

	return fmt.Errorf("%s: not found", cmd)
}

func pwdHandler(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Println(dir)
	return nil
}

func cdHandler(args []string) error {
	relPath := args[0]
	absolutePath := utils.ExpandPath(relPath)

	if err := os.Chdir(absolutePath); err != nil {
		return fmt.Errorf("cd: %s: No such file or directory", relPath)
	}
	return nil
}

func getHandler(cmd string) (handler, error) {
	if builtinHandler, exists := builtins[cmd]; exists {
		return builtinHandler, nil
	}

	if absPath, err := utils.GetFromPath(cmd); err == nil {
		ex := exec.Command(absPath)
		ex.Args[0] = cmd
		ex.Stdout = os.Stdout
		ex.Stdin = os.Stdin
		ex.Stderr = os.Stderr

		return func(args []string) error {
			ex.Args = append(ex.Args, args...)
			return ex.Run()
		}, nil
	}

	return nil, fmt.Errorf("%s: command not found", cmd)
}

var builtins = make(map[string]handler)

func init() {
	builtins["exit"] = exitHandler
	builtins["echo"] = echoHandler
	builtins["type"] = typeHandler
	builtins["pwd"] = pwdHandler
	builtins["cd"] = cdHandler
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("$ ")
		if !scanner.Scan() {
			break
		}

		fields := strings.Fields(scanner.Text())
		cmd, args := fields[0], fields[1:]

		handler, err := getHandler(cmd)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = handler(args)
		if err == errExit {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}
