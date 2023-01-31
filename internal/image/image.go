package image

import (
	"fmt"
	"github.com/mholt/archiver/v3"
	"os"
	"os/exec"
)

var basePath = "/home/nicolas/projects/havnearbeider"
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
	// Path is where the archive is located
	Path string
}

func ArchivedImageFromName(name string) ArchivedImage {
	return ArchivedImage{
		Name: name,
		Path: fmt.Sprintf("%s/%s.tar", imageBasePath, name),
	}
}

func (i ArchivedImage) Load() LoadedImage {
	extractPath := fmt.Sprintf("%s/%s", extractBasePath, i.Name)
	return i.LoadToPath(extractPath)
}

func (i ArchivedImage) LoadToPath(path string) LoadedImage {
	must(archiver.Unarchive(i.Path, path))

	return LoadedImage{
		Name: i.Name,
		Root: path,
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

func (i LoadedImage) Package() ArchivedImage {
	return i.PackageToPath(imageBasePath)
}

func (i LoadedImage) PackageToPath(path string) ArchivedImage {
	files, err := os.ReadDir(i.Root)
	must(err)

	toArchive := []string{}
	for _, f := range files {
		toArchive = append(toArchive, fmt.Sprintf("%s/%s", i.Root, f.Name()))
	}

	tarPath := fmt.Sprintf("%s/%s.tar", path, i.Name)
	must(archiver.Archive(toArchive, tarPath))

	return ArchivedImage{
		Name: i.Name,
		Path: tarPath,
	}
}
