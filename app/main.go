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
		op := parser.OpNil
		for _, cmd := range cmds {
			proc := &handlers.Proc{Args: cmd.Args}
			if op != parser.OpNil {
				proc.Stdin = buf
			} else {
				proc.Stdin = os.Stdin
			}

			switch cmd.Next {
			case parser.OpRedirectStdOut:
				proc.Stdout = buf
				proc.Stderr = os.Stderr
			case parser.OpRedirectStdErr:
				proc.Stdout = os.Stdout
				proc.Stderr = buf
			case parser.OpRedirectBoth:
				proc.Stdout = buf
				proc.Stderr = buf
			default:
				proc.Stdout = os.Stdout
				proc.Stderr = os.Stderr
			}

			var handler handlers.Handler
			switch op {
			case parser.OpRedirectStdOut, parser.OpRedirectStdErr, parser.OpRedirectBoth:
				handler, err = handlers.GetRedirectHandler(cmd.Name, op)
			default:
				handler, err = handlers.GetHandler(cmd.Name)
			}

			if err != nil {
				if errors.Is(err, handlers.ErrCommandNotFound) {
					fmt.Println(err)
					continue
				}

				fmt.Println(err)
			}

			err = handler(proc)
			if err == handlers.ErrExit {
				break outer
			}
			if cmd.Next == parser.OpNil {
				buf.Reset()
			}

			op = cmd.Next
		}
	}
}
