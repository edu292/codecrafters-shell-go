package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
		switch cmd {
		case "exit":
			break outer
		case "echo":
			fmt.Println(strings.Join(args, " "))
		default:
			fmt.Printf("%s: command not found\n", cmd)
		}
	}
}
