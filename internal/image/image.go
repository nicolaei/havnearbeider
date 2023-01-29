package image

import (
	"fmt"
	"github.com/mholt/archiver/v3"
	"os/exec"
)

var basePath = "/home/nicolas/projects/container-runtime"
var imageBasePath = fmt.Sprintf("%s/images", basePath)
var extractBasePath = fmt.Sprintf("%s/filesystems", basePath)

func must(err error) {
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			return
		default:
			panic(err)
		}
	}
}

type ArchivedImage struct {
	// Name is the name of this image
	Name string
}

func (i ArchivedImage) Load() LoadedImage {
	imagePath := fmt.Sprintf("%s/%s.tar", imageBasePath, i.Name)
	extractPath := fmt.Sprintf("%s/%s", extractBasePath, i.Name)

	must(archiver.Unarchive(imagePath, extractPath))

	return LoadedImage{
		Name: i.Name,
		Root: extractPath,
	}
}

type LoadedImage struct {
	// Name is the name of this image
	Name string
	// Root is the root path of the filesystem
	Root string
}

func LoadedImageFromName(name string) LoadedImage {
	return LoadedImage{
		Name: name,
		Root: fmt.Sprintf("%s/%s", extractBasePath, name),
	}
}
