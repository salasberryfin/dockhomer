package container

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/salasberryfin/dockhomer/utils"
)

// ImageManifest describes the information contained in `manifest.json` for each image
type ImageManifest struct {
	Config   string
	RepoTags []string
	Layers   []string
}

// Image holds basic information about existing filesystems stored in the host machine
type Image struct {
	Name string
	Hex  string
	Path string
}

// ListImages retrieves the existing images in the host machine
func ListImages() ([]Image, error) {
	if _, err := os.Stat(dockhomerImagesRoot); err != nil {
		if os.IsNotExist(err) {
			log.Printf("%s does not exist, creating...\n", dockhomerImagesRoot)
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

	return images, nil
}

// pullImage checks if the image exists in the default image folder
// if not, it downloads the filesystem to the root folder
// pending a better implementation
// for now, it only supports single-layer images: need to check what happens when using
// a multi-layer image
func pullImage(name string, force bool) (Image, error) {
	// if force is set to false, check if image already exists in the root dockhomer directory for images
	if !force {
		img, err := InspectImage(name)
		if err != nil {
			// image is not in the images path -> pull
			log.Printf("Image not found in %s\n", dockhomerImagesRoot)
		} else {
			log.Printf("Image %s already exists in %s - Skipping pull.\n", name, dockhomerImagesRoot)
			return img, nil
		}
	}

	log.Printf("Pulling %s\n", name)
	img, err := crane.Pull(name)
	manifest, err := img.Manifest()
	manifestLayer := manifest.Layers[0].Digest.Hex
	if err != nil {
		return Image{}, err
	}
	imageRootPath := filepath.Join(dockhomerImagesRoot, name)
	err = createDirs([]string{imageRootPath}, fs.FileMode(0755))
	if err != nil {
		return Image{}, err
	}
	tarballPath := filepath.Join(imageRootPath, name+".tar")
	err = crane.Save(img, "dockhomervers", tarballPath)
	if err != nil {
		return Image{}, err
	}
	r, err := os.Open(tarballPath)
	if err != nil {
		return Image{}, err
	}
	defer r.Close()
	err = utils.Untar(r, imageRootPath)
	image := Image{
		Name: name,
		Hex:  manifestLayer,
		Path: imageRootPath,
	}
	log.Printf("Pulled image %v", image)

	return image, nil
}

// InspectImage returns the relevant information about the specified image if it
// is found in the root image directory
func InspectImage(name string) (Image, error) {
	var image Image
	imgs, err := ListImages()
	if err != nil {
		return Image{}, err
	}

	for _, img := range imgs {
		if img.Name == name {
			manifestJsonPath := filepath.Join(img.Path, "manifest.json")
			content, err := ioutil.ReadFile(manifestJsonPath)
			if err != nil {
				return image, err
			}
			var manifestJson []ImageManifest
			err = json.Unmarshal(content, &manifestJson)
			if err != nil {
				return image, err
			}
			image.Name = name
			image.Hex = strings.Split(manifestJson[0].Layers[0], ".")[0]
			image.Path = img.Path
			fmt.Printf("The image is %v\n", image)

			return image, nil
		}
	}

	return Image{}, fmt.Errorf("Image '%s' not found in '%s'\n", name, dockhomerImagesRoot)
}
