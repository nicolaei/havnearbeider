package main

import (
	"fmt"
	"github.com/nicolaei/container-runtime/internal/runtime"
	"os"
)

func main() {
	switch os.Args[1] {
	case "__run__":
		container := runtime.Container{
			Image:   os.Args[2],
			Command: os.Args[3:],
		}

		container.Run()
	case "kjør":
		container := runtime.Container{
			Image:   os.Args[2],
			Command: os.Args[3:],
		}
		container.Create()
	case "bygg":
		buildContainer()
	default:
		panic("invalid argument")
	}
}

// buildContainer builds a container from a .containerspec file
func buildContainer() {
	fmt.Printf("building: %s\n", os.Args[2])
}
