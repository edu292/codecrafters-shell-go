package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Print("$ ")
	var command string
	_, err := fmt.Scanf("%s", &command)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("%s: command not found", command)
}
