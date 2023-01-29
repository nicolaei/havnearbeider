package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type Container struct {
	Image   string
	Command []string
}

// Run runs the specified command on the given image.
func (c Container) Run() {
	fmt.Printf("running: %s on %s as PID %d\n", os.Args[3:], os.Args[2], os.Getpid())

	filesystem := getFilesystem(os.Args[2])
	unmount := mountFilesystem(filesystem)
	defer unmount()

	exit_func := newProcessSpace()
	defer exit_func()

	newHostname()

	runCommand()
}

// Create creates the container by creating a new namespace and attaching CGroups.
func (c Container) Create() {
	cmd := exec.Command("/proc/self/exe", append([]string{"__run__"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}
