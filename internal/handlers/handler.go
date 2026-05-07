package handlers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"codecrafters-shell-go/internal/parser"
	"codecrafters-shell-go/internal/utils"
)

type Proc struct {
	Args   []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type Handler func(p *Proc) error

var (
	ErrExit            = errors.New("exit")
	ErrCommandNotFound = errors.New("command not found")
)

func exitHandler(p *Proc) error {
	return ErrExit
}

func echoHandler(p *Proc) error {
	fmt.Fprintln(p.Stdout, strings.Join(p.Args, " "))

	return nil
}

func typeHandler(p *Proc) error {
	if len(p.Args) == 0 {
		return nil
	}

	name := p.Args[0]
	if _, exists := BuiltIns[name]; exists {
		fmt.Fprintf(p.Stdout, "%s is a shell builtin\n", name)
		return nil
	}

	if absPath, err := utils.GetFromPath(name); err == nil {
		fmt.Fprintf(p.Stdout, "%s is %s\n", name, absPath)
		return nil
	}

	fmt.Fprintf(p.Stderr, "%s: not found\n", name)

	return nil
}

func pwdHandler(p *Proc) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Fprintln(p.Stdout, dir)
	return nil
}

func cdHandler(p *Proc) error {
	var relPath string
	if len(p.Args) == 0 {
		relPath, _ = os.UserHomeDir()
	} else {
		relPath = p.Args[0]
	}

	absPath := utils.ExpandPath(relPath)

	if err := os.Chdir(absPath); err != nil {
		fmt.Fprintf(p.Stderr, "cd: %s: No such file or directory\n", relPath)
	}
	return nil
}

func GetHandler(name string) (Handler, error) {
	if name == "" {
		goto error
	}

	if builtinHandler, exists := BuiltIns[name]; exists {
		return builtinHandler, nil
	}

	if absPath, err := utils.GetFromPath(name); err == nil {
		ex := exec.Command(absPath)
		ex.Args[0] = name

		return func(p *Proc) error {
			ex.Args = append(ex.Args, p.Args...)
			ex.Stdin = p.Stdin
			ex.Stdout = p.Stdout
			ex.Stderr = p.Stderr
			return ex.Run()
		}, nil
	}

error:
	return nil, fmt.Errorf("%s: %w", name, ErrCommandNotFound)
}

func GetRedirectHandler(relPath string, op parser.Op) (Handler, error) {
	return func(p *Proc) error {
		f, err := os.OpenFile(utils.ExpandPath(relPath), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, p.Stdin)
		return err
	}, nil
}

var BuiltIns = make(map[string]Handler)

func init() {
	BuiltIns["exit"] = exitHandler
	BuiltIns["echo"] = echoHandler
	BuiltIns["type"] = typeHandler
	BuiltIns["pwd"] = pwdHandler
	BuiltIns["cd"] = cdHandler
}
