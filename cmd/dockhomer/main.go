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
		cont := container.New("/tmp/dockhomer")
		log.Printf("Creating container with id: %d\n", cont.ID)
		if command == "shell" {
			//container.New("/tmp/dockhomer/")
			cont.OpenShell()
		} else {
			container.Run()
		}
	default:
		log.Fatalf("No valid command found\n")
	}
}
