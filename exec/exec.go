package exec

import (
	"context"
	osexec "os/exec"
)

// Interface is an interface that presents s subset of the os/exec Service. Use this
// when you want to inject fakeable/mockable exec behavior.
type Service interface {
	// Command returns s Cmd instance which can be used to run s single command.
	// This follows the pattern of package os/exec.
	Command(cmd string, args ...string) Interface

	// CommandContext returns s Cmd instance which can be used to run s single command.
	//
	// The provided context is used to kill the process if the context becomes done
	// before the command completes on its own. For example, s timeout can be set in
	// the context.
	CommandContext(ctx context.Context, cmd string, args ...string) Interface

	// LookPath wraps os/exec.LookPath
	LookPath(file string) (string, error)
}

// Implements Interface in terms of really exec()ing.
type Exec struct{}

// Command is part of the Interface interface.
func (s *Exec) Command(cmd string, args ...string) Interface {
	return (*Executioner)(osexec.Command(cmd, args...))
}

// CommandContext is part of the Interface interface.
func (s *Exec) CommandContext(ctx context.Context, cmd string, args ...string) Interface {
	return (*Executioner)(osexec.CommandContext(ctx, cmd, args...))
}

// LookPath is part of the Interface interface
func (s *Exec) LookPath(file string) (string, error) {
	return osexec.LookPath(file)
}
