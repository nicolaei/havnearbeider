package build

import (
	"fmt"
	"github.com/nicolaei/havnearbeider/internal/image"
	"github.com/nicolaei/havnearbeider/internal/runtime"
	"strings"
)

type Step interface {
	// Description is used to tell the user what the step is doing
	Description() string
	// Exec executes the step and stores the result in outDirectory.
	Exec(workingDirectory string, baseImage *image.ArchivedImage) image.LoadedImage
}

func StepFromLine(line string) (Step, error) {
	splitLine := strings.Split(line, " ")

	switch splitLine[0] {
	case "FRA":
		return ImageStep{splitLine[1]}, nil
	case "KJØR":
		return RunStep{splitLine[1:]}, nil
	default:
		return nil, fmt.Errorf("ukjent direktiv %s for linje %s\n", splitLine[0], splitLine)
	}
}

// ImageStep takes in a name of an image we want to base of, and gives us the filesystem.
type ImageStep struct {
	ImageName string
}

func (s ImageStep) Description() string {
	return fmt.Sprintf("getting image for %s", s.ImageName)
}

func (s ImageStep) Exec(outDirectory string, baseImage *image.ArchivedImage) image.LoadedImage {
	toLoad := image.ArchivedImageFromName(s.ImageName)

	return toLoad.LoadToPath(outDirectory)
}

// RunStep is a step that modifies the image by running a command
type RunStep struct {
	Command []string
}

func (s RunStep) Description() string {
	return fmt.Sprintf("running %s", s.Command)
}

func (s RunStep) Exec(outDirectory string, baseImage *image.ArchivedImage) image.LoadedImage {
	if baseImage == nil {
		panic("step requires a base image")
	}

	buildEnvironment := runtime.Container{
		Image:   *baseImage,
		Command: s.Command,
	}

	return buildEnvironment.Create()
}
