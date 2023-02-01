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
		id, err := container.Create()
		if err != nil {
			log.Fatalf("Failed to create container: %s\n", err)
		}
		if command == "shell" {
			log.Printf("An interactive shell will be open for container %d\n", id)
			//container.RunShell()
		} else {
			container.Run(os.Args[2:])
		}
	default:
		log.Fatalf("No valid command found\n")
	}
}
