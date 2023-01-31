package main

import (
	"github.com/nicolaei/havnearbeider/internal/build"
	"github.com/nicolaei/havnearbeider/internal/image"
	"github.com/nicolaei/havnearbeider/internal/runtime"
	"os"
)

func main() {
	switch os.Args[1] {
	case "__run__":
		container := runtime.RunningContainer{
			Image:   image.LoadedImageFromName(os.Args[2]),
			Command: os.Args[3:],
		}

		container.Run()
	case "kjør":
		container := runtime.Container{
			Image:   image.ArchivedImageFromName(os.Args[2]),
			Command: os.Args[3:],
		}
		container.Create()
	case "bygg":
		buildSpec := build.NewSpecFromFile(os.Args[2], os.Args[3])

		buildSpec.Build()
	default:
		panic("invalid argument")
	}
}
