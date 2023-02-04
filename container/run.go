package container

import (
	"log"
	"os"
	"os/exec"
)

const (
	dockhomerHome    = "/tmp/dockhomer"
	dockhomerVolumes = dockhomerHome + "/volumes"
	defaultShell     = "/bin/sh"
)

// RunShell opens an interactive bash shell
func RunShell() {
	hostname := os.Args[1] // hostname defaults to container ID

	err := newFilesystem(dockhomerHome, hostname)
	if err != nil {
		log.Fatalf("Failed to build file system: %v\n", err)
	}

	cmd := exec.Command(defaultShell)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run the command: %s\n", err)
	}
}

// Run runs a command in a new container
func Run() {
	hostname := os.Args[1]
	term := os.Args[2:]

	err := newFilesystem(dockhomerHome, hostname)
	if err != nil {
		log.Fatalf("Failed to build file system: %v\n", err)
	}

	cmd := exec.Command(term[0], term[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
