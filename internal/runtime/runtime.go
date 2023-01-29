package runtime

import (
	"fmt"
	"github.com/nicolaei/container-runtime/internal/image"
	"os"
	"os/exec"
	"syscall"
)

type RunningContainer struct {
	Image   image.LoadedImage
	Command []string
}

// Run runs the specified command on the given image.
func (c RunningContainer) Run() {
	fmt.Printf("running: %s on %s\n", c.Command, c.Image.Name)

	unmount := mountFilesystem(c.Image.Root)
	defer unmount()

	exit_func := newProcessSpace()
	defer exit_func()

	newHostname()

	runCommand(c.Command)
}

type Container struct {
	Image   image.ArchivedImage
	Command []string
}

// Create creates the container by creating a new namespace and attaching CGroups.
func (c Container) Create() {
	c.Image.Load()

	cmd := exec.Command("/proc/self/exe", append([]string{"__run__", c.Image.Name}, c.Command...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}
