package container

import (
	"log"
	"os/exec"
)

// Interactive opens an interactive bash shell
func Interactive() {
	cmd := exec.Command("/bin/bash")
	create(cmd) // create new container
	log.Printf("Interactive shell is loading %v\n", cmd.Args)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// Run runs a command in a new container
func Run(term []string) {
	cmd := exec.Command(term[0], term[1])
	create(cmd)	// create new container
	log.Printf("Running %s in the container\n", term)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
