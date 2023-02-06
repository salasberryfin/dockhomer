package container

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

// Image holds basic information about existing filesystems stored in the host machine
type Image struct {
	Name string
	Path string
}

const (
	dockhomerRoot           = "/tmp/dockhomer"
	dockhomerImagesRoot     = dockhomerRoot + "/images"
	dockhomerContainersRoot = dockhomerRoot + "/containers"
)

// ListImages retrieves the existing images in the host machine
func ListImages() ([]Image, error) {
	if _, err := os.Stat(dockhomerImagesRoot); err != nil {
		if os.IsNotExist(err) {
			log.Printf("%s does not exists, creating...\n", dockhomerImagesRoot)
			err := createDirs([]string{dockhomerImagesRoot}, fs.FileMode(0755))
			if err != nil {
				return nil, err
			}
		}
	}
	entries, err := os.ReadDir(dockhomerImagesRoot)
	if err != nil {
		return nil, err
	}
	images := make([]Image, 0, len(entries))
	for _, entry := range entries {
		image := Image{
			Name: entry.Name(),
			Path: filepath.Join(dockhomerImagesRoot, entry.Name()),
		}
		images = append(images, image)
	}

	log.Printf("Images:\n%v\n", images)

	return images, nil
}

// pullImage checks if the image exists in the default image folder
// the current behavior is to use a hardcoded `./volumes` folder in the project
// root in which image filesystems should be stored
// pending a better implementation
func pullImage(name string) error {

	return nil
}

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

// createDirs ships the container with the required directories
// it defaults to using a temporary folder ./volumes for any container dir
func createDirs(dirs []string, perm fs.FileMode) error {
	for _, dir := range dirs {
		log.Printf("Directories will be created in %s\n", dir)
		err := os.MkdirAll(dir, perm)
		if err != nil {
			return err
		}
	}

	return nil
}
