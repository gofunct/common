package common

import (
	"context"
	"github.com/gofunct/common/exec"
)

// Interface is an interface that presents a subset of the os/exec API. Use this
// when you want to inject fakeable/mockable exec behavior.
type Exec interface {
	// Command returns a Cmd instance which can be used to run a single command.
	// This follows the pattern of package os/exec.
	Command(cmd string, args ...string) exec.Interface

	// CommandContext returns a Cmd instance which can be used to run a single command.
	//
	// The provided context is used to kill the process if the context becomes done
	// before the command completes on its own. For example, a timeout can be set in
	// the context.
	CommandContext(ctx context.Context, cmd string, args ...string) exec.Interface

	// LookPath wraps os/exec.LookPath
	LookPath(file string) (string, error)
}

// New returns a new Interface which will os/exec to run commands.
func NewExec() Exec {
	return &exec.API{}
}
