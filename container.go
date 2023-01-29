package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "__child__":
		runContainer()
	case "kjør":
		createContainer()
	case "bygg":
		buildContainer()
	default:
		panic("invalid argument")
	}
}

// createContainer sets up the environment where our container will run.
func createContainer() {
	cmd := exec.Command("/proc/self/exe", append([]string{"__child__"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

// runContainer runs the specified container
func runContainer() {
	fmt.Printf("running: %s on %s as PID %d\n", os.Args[3:], os.Args[2], os.Getpid())

	filesystem := getFilesystem(os.Args[2])
	unmount := mountFilesystem(filesystem)
	defer unmount()

	exit_func := newProcessSpace()
	defer exit_func()

	newHostname()

	runCommand()
}

// buildContainer builds a container from a .containerspec file
func buildContainer() {
	fmt.Printf("building: %s\n", os.Args[2])
}
