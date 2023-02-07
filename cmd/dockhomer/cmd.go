package main

import (
	"log"
	"os"

	"github.com/salasberryfin/dockhomer/container"
	"github.com/spf13/cobra"
)

type DockhomerOptions struct {
	Args []string
	//ConfigFlags
}

// NewDefaultCommand generates the `dockhomer` command with default arguments
func NewDefaultCommand() *cobra.Command {
	return DockhomerOptions{
		Args: os.Args,
	}
}

func runHelp(cmd *cobra.Command) {
	cmd.Help()
}

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
