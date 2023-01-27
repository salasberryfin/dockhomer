package main

import (
	"log"
	"os"

	"github.com/salasberryfin/dockhomer/container"
)

func main() {
	// read user argument
	if len(os.Args) > 3 {
		log.Fatalf("Wrong number of arguments\n")
	}
	action, command := os.Args[1], os.Args[2]
	switch action {
	case "run":
		if command == "shell" {
			container.Interactive()
		}
	default:
		log.Fatalf("No valid command found\n")
	}
}
