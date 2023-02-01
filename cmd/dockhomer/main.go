package main

import (
	"log"
	"os"

	"github.com/salasberryfin/dockhomer/container"
)

func main() {
	action, command := os.Args[1], os.Args[2]
	switch action {
	case "run":
		if command == "shell" {
			container.RunShell()
		} else {
			container.Run(os.Args[2:])
		}
	default:
		log.Fatalf("No valid command found\n")
	}
}
