// +build !linux

package nsenter

import (
	"context"
	"fmt"

	"k8s.io/utils/exec"
)

const (
	// DefaultHostRootFsPath is path to host's filesystem mounted into container
	// with kubelet.
	DefaultHostRootFsPath = "/rootfs"
)

// Nsenter is a type alias for backward compatibility
type Nsenter = NSEnter

// NSEnter is part of experimental support for running the kubelet
// in a container.
type NSEnter struct {
	// a map of commands to their paths on the host filesystem
	Paths map[string]string
}

// NewNsenter constructs a new instance of NSEnter
func NewNsenter(hostRootFsPath string, executor exec.Interface) (*Nsenter, error) {
	return &Nsenter{}, nil
}

// Exec executes nsenter commands in hostProcMountNsPath mount namespace
func (ne *NSEnter) Exec(cmd string, args []string) exec.Cmd {
	return nil
}

// AbsHostPath returns the absolute runnable path for a specified command
func (ne *NSEnter) AbsHostPath(command string) string {
	return ""
}

// SupportsSystemd checks whether command systemd-run exists
func (ne *NSEnter) SupportsSystemd() (string, bool) {
	return "", false
}

// Command returns a command wrapped with nenter
func (ne *NSEnter) Command(cmd string, args ...string) exec.Cmd {
	return nil
}

// CommandContext returns a CommandContext wrapped with nsenter
func (ne *NSEnter) CommandContext(ctx context.Context, cmd string, args ...string) exec.Cmd {
	return nil
}

// LookPath returns a LookPath wrapped with nsenter
func (ne *NSEnter) LookPath(file string) (string, error) {
	return "", fmt.Errorf("not implemented, error looking up : %s", file)
}

var _ exec.Interface = &NSEnter{}
