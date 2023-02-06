package container

import (
	"fmt"
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
	Image    string
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

// OpenShell opens an interactive shell in the container
func (c *Container) OpenShell() error {
	cmd := reexec.Command("shell", strconv.Itoa(c.ID))
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
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

// RunCmd runs the command passed as CLI argument in the container
func (c *Container) RunCmd(args ...string) error {
	newArgs := []string{"command", strconv.Itoa(c.ID)}
	newArgs = append(newArgs, args...)
	cmd := reexec.Command(newArgs...)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
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
func New(image, root string) *Container {
	// network configuration test
	//network.NewBridge("dockhomer-bridge")
	//network.NewVethPair("vethdockhomer0")
	_, err := ListImages()
	if err != nil {
		log.Fatalf("Failed to list existing images: %v\n", err)
	}

	return &Container{
		ProcAttr: syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS |
				syscall.CLONE_NEWPID |
				syscall.CLONE_NEWUSER |
				syscall.CLONE_NEWNS |
				syscall.CLONE_NEWNET,
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
