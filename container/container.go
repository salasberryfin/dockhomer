package container

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

// generateID creates a random identifier
func generateID() int {
	r := make([]byte, 4)
	rand.Seed(time.Now().UnixNano())
	_, err := rand.Read(r)
	if err != nil {
		log.Fatal(err)
	}
	id, err := strconv.Atoi(fmt.Sprintf("%d%d%d%d", r[0], r[1], r[2], r[3]))

	return id
}

// create configures the namespaces for a new container
func create(cmd *exec.Cmd) {
	containerID := generateID()
	log.Printf("Starting a new container %d\n", containerID)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 1,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 1,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}
}
