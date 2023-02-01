package container

import (
	"log"
	"os"
	"os/exec"
)

// RunShell opens an interactive bash shell
func RunShell() {
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("Interactive shell is loading %v\n", cmd.Args)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run the command: %s\n", err)
	}
}

// Run runs a command in a new container
func Run(term []string) {
	cmd := exec.Command(term[0], term[1:]...)
	_, err := Create() // create new container
	if err != nil {
		log.Fatalf("Failed to create container: %s\n", err)
	}
	log.Printf("Running %s in the container\n", term)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
