package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"

	"codecrafters-shell-go/internal/handlers"
	"codecrafters-shell-go/internal/parser"
)

type state int

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	buf := new(bytes.Buffer)
	var err error

outer:
	for {
		fmt.Print("$ ")
		if !scanner.Scan() {
			break
		}

		cmds := parser.ParseInput(scanner.Bytes())
		op := parser.Nil
		for _, cmd := range cmds {
			proc := &handlers.Proc{Args: cmd.Args}
			if op != parser.Nil {
				proc.Stdin = buf
			} else {
				proc.Stdin = os.Stdin
			}

			switch {
			case parser.Is(cmd.Op, parser.Stdout):
				proc.Stdout = buf
				proc.Stderr = os.Stderr
			case parser.Is(cmd.Op, parser.Stderr):
				proc.Stdout = os.Stdout
				proc.Stderr = buf
			case parser.Is(cmd.Op, parser.Both):
				proc.Stdout = buf
				proc.Stderr = buf
			default:
				proc.Stdout = os.Stdout
				proc.Stderr = os.Stderr
			}

			var handler handlers.Handler
			if parser.Is(op, parser.Redir) {
				handler, err = handlers.GetRedirectHandler(cmd.Name, parser.Is(op, parser.Append))
			} else {
				handler, err = handlers.GetHandler(cmd.Name)
			}

			if errors.Is(err, handlers.ErrCommandNotFound) {
				fmt.Println(err)
				continue
			}

			err = handler(proc)
			if err != nil {
				if err == handlers.ErrExit {
					break outer
				}
				fmt.Println(err)
			}

			if cmd.Op == parser.Nil {
				buf.Reset()
			}

			op = cmd.Op
		}
	}
}
