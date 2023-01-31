package container

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

const (
	dockhomerHome    = "/tmp/dockhomer"
	dockhomerVolumes = dockhomerHome + "/volumes"
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
func create(cmd *exec.Cmd) (int, error) {
	containerID := generateID()
	log.Printf("Starting a new container %d\n", containerID)
	err := createDirs([]string{"/proc"})
	if err != nil {
		return 0, err
	}
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
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	return containerID, nil
}

// createDirs ships the container with the required directories
// it defaults to using a temporary folder ./volumes for any container dir
func createDirs(dirs []string) error {
	for _, dir := range dirs {
		dirPath := filepath.Join(dockhomerVolumes, dir)
		log.Printf("Container directories will be created in %s\n", dirPath)
		err := os.MkdirAll(dirPath, 0750)
		if err != nil {
			return err
		}
	}

	return nil
}
