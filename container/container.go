package container

import (
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/reexec"
)

type Container struct {
	ProcAttr syscall.SysProcAttr
	ID       int
	Root     string
	Stdin    *os.File
	Stdout   *os.File
	Stderr   *os.File
}

func init() {
	reexec.Register("shell", RunShell)
	reexec.Register("command", Run)
	if reexec.Init() {
		os.Exit(0)
	}
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

func (c *Container) OpenShell() error {
	cmd := reexec.Command("shell", strconv.Itoa(c.ID))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &c.ProcAttr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to initialize container namespace: %v\n", err)
		return err
	}
	if err := cmd.Wait(); err != nil {
		log.Fatalf("Failed to start container: %v\n", err)
		return err
	}

	return nil
}

// New creates a new instance of a container
func New(root string) *Container {
	return &Container{
		ProcAttr: syscall.SysProcAttr{
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
		},
		ID:     generateID(),
		Root:   root,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
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
