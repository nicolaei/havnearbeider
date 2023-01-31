package build

import (
	"bufio"
	"fmt"
	"github.com/nicolaei/havnearbeider/internal/image"
	"log"
	"os"
)

// TODO: Use a JSON file for the command to run

type Spec struct {
	Name  string
	Steps []Step
}

func (s Spec) Build() image.ArchivedImage {
	fmt.Printf("[...] Building %s\n", s.Name)

	builtLayer := s.buildStep(s.Steps[0], nil, 1, len(s.Steps))

	for index, step := range s.Steps[1:] {
		builtLayer = s.buildStep(step, &builtLayer, index+2, len(s.Steps))
	}

	fmt.Printf("[...] Packaging final image...\n")

	// TODO: This could probably be redone better by having a `CopyArchiveImage` function
	finalImage := builtLayer.LoadToPath("build/final")
	return finalImage.Package()
}

func (s Spec) buildStep(step Step, baseLayer *image.ArchivedImage, curretStep int, totalSteps int) image.ArchivedImage {
	fmt.Printf("[%d/%d] %s\n", curretStep, totalSteps, step.Description())

	stepBase := fmt.Sprintf("build/step-%d", curretStep)

	builtLayer := step.Exec(stepBase, baseLayer)
	builtLayer.Name = fmt.Sprintf("step-%d", curretStep)

	return builtLayer.PackageToPath(stepBase)
}

// NewSpecFromFile reads a container spec file and parses it to create a Spec
func NewSpecFromFile(file string, imageName string) Spec {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	scanner := bufio.NewScanner(f)

	steps := []Step{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		step, err := StepFromLine(line)
		if err != nil {
			log.Fatal(err)
		}

		steps = append(steps, step)
	}

	return Spec{
		Name:  imageName,
		Steps: steps,
	}
}
