package main

import (
	"fmt"
	"log"
)

func main() {
	for {
		fmt.Print("$ ")
		var command string
		_, err := fmt.Scanf("%s", &command)
		if err != nil {
			log.Fatalf("%v", err)
		}
		if command == "exit" {
			break
		}

		fmt.Printf("%s: command not found\n", command)
	}
}
