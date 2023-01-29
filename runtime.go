package main

import (
	"fmt"
	"github.com/mholt/archiver/v3"
	"os"
	"os/exec"
	"syscall"
)

func getFilesystem(image string) string {
	basePath := "/home/nicolas/projects/container-runtime/filesystems"

	imageBase := fmt.Sprintf("%s/%s", basePath, image)
	must(archiver.Unarchive(fmt.Sprintf("%s.tar", imageBase), imageBase))

	return imageBase
}

func mountFilesystem(path string) func() {
	exit, err := chroot(path)

	must(err)

	must(os.Chdir("/"))

	return func() { must(exit()) }
}

func newHostname() {
	err := syscall.Sethostname([]byte("container"))
	if err != nil {
		return
	}
}

func newProcessSpace() func() {
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	return func() { must(syscall.Unmount("proc", 0)) }
}

func runCommand() {
	cmd := exec.Command(os.Args[3], os.Args[4:]...)
	cmd.Env = append(cmd.Env, "PATH=/bin:/usr/bin")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
}

/// ----- Boilerplate from here ----- ///

func chroot(path string) (func() error, error) {
	root, err := os.Open("/")
	if err != nil {
		return nil, err
	}

	if err := syscall.Chroot(path); err != nil {
		root.Close()
		return nil, err
	}

	return func() error {
		defer root.Close()
		if err := root.Chdir(); err != nil {
			return err
		}
		return syscall.Chroot(".")
	}, nil
}

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
