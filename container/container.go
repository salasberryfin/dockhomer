package container

import (
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/reexec"
)

const (
	dockhomerHome    = "/tmp/dockhomer"
	dockhomerVolumes = dockhomerHome + "/volumes"
)

func init() {
	reexec.Register("namespaceConf", namespaceConf)
	if reexec.Init() {
		os.Exit(0)
	}
}

// namespaceConf is used by reexec for namespace initialization, including
// pivot_root to define the container root filesystem
func namespaceConf() {
	newRootFs := os.Args[1]
	target := filepath.Join(newRootFs, "proc")

	if _, err := os.Stat(target); os.IsNotExist(err) {
		log.Printf("%s does not exist, creating...\n", target)
		if err := createDirs([]string{target}, fs.FileMode(0755)); err != nil {
			log.Fatalf("Failed to create directory: %s\n", err)
			os.Exit(1)
		}
	}
	log.Printf("Mounting %s\n", target)
	if err := syscall.Mount("proc", target, "proc", 0, ""); err != nil {
		log.Fatalf("Failed to mount %s: %s\n", target, err)
		os.Exit(1)
	}

	if err := pivotRoot(newRootFs, "oldFs"); err != nil {
		log.Fatal("Failed to pivot_root: %s\n", err)
		os.Exit(1)
	}
	RunShell()
}

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

// Create configures the namespaces for a new container
func Create() (int, error) {
	containerID := generateID()
	log.Printf("Starting a new container %d\n", containerID)

	cmd := reexec.Command("namespaceConf", dockhomerHome)

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
		Credential: &syscall.Credential{
			Uid: 0,
			Gid: 0,
		},
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to initialize container namespace\n")
		os.Exit(1)
	}
	if err := cmd.Wait(); err != nil {
		log.Printf("Failed to start container\n")
		os.Exit(1)
	}

	return containerID, nil
}

// createDirs ships the container with the required directories
// it defaults to using a temporary folder ./volumes for any container dir
func createDirs(dirs []string, perm fs.FileMode) error {
	for _, dir := range dirs {
		//dirPath := filepath.Join(dockhomerVolumes, dir)
		log.Printf("Container directories will be created in %s\n", dir)
		err := os.MkdirAll(dir, perm)
		if err != nil {
			return err
		}
	}

	return nil
}
