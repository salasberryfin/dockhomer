package container

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

// newFileSystem mounts all the required directories to create an isolated
// file system for the new container. This allows using different system images
func newFilesystem(newRootFs, hostname string) error {
	target := filepath.Join(newRootFs, "proc")

	if _, err := os.Stat(target); os.IsNotExist(err) {
		log.Printf("%s does not exist, creating...\n", target)
		if err := createDirs([]string{target}, fs.FileMode(0755)); err != nil {
			return err
		}
	}

	log.Printf("Mounting new /proc in %s\n", target)
	if err := syscall.Mount("proc", target, "proc", 0, ""); err != nil {
		return err
	}

	log.Printf("pivot_root to new filesystem root %s\n", newRootFs)
	if err := pivotRoot(newRootFs, "oldFs"); err != nil {
		return err
	}

	log.Printf("Setting new hostname: %s\n", hostname)
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		return err
	}

	return nil
}

// pivotRoot allows to set a new root filesystem for the container
func pivotRoot(newFs, oldFolder string) error {
	log.Printf("Pivoting root filesystem to %s\n", newFs)
	err := syscall.Mount(newFs, newFs, "", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		return err
	}

	oldFs := filepath.Join(newFs, oldFolder)
	if err := createDirs([]string{oldFs}, fs.FileMode(0700)); err != nil {
		return err
	}

	if err := syscall.PivotRoot(newFs, oldFs); err != nil {
		return err
	}

	if err := os.Chdir("/"); err != nil {
		return err
	}

	if err := syscall.Unmount(oldFolder, syscall.MNT_DETACH); err != nil {
		return err
	}

	if err := os.RemoveAll(oldFolder); err != nil {
		return err
	}

	return nil
}

// mountFileSystem uses syscall.mount the file system to the
// mount namespace of the container
func mountFileSystem() error {

	return nil
}

// unmountFileSystem uses syscall.unmount the file system to the
// unmount namespace of the container
func unmountFileSystem() error {

	return nil
}
