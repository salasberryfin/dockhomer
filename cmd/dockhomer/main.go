package main

import (
	"log"
	"os"

	"github.com/salasberryfin/dockhomer/container"
)

func main() {
	action, image, command := os.Args[1], os.Args[2], os.Args[3]
	switch action {
	case "run":
		cont := container.New(image, "/tmp/dockhomer")
		log.Printf("Creating container with id: %d\n", cont.ID)
		if command == "shell" {
			cont.OpenShell()
		} else {
			cont.RunCmd(os.Args[3:]...)
		}
	default:
		log.Fatalf("No valid command found\n")
	}
}
