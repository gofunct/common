package exec

import (
	"context"
	"io"
	osexec "os/exec"
)

// Cmd is an interface that presents an API that is very similar to Cmd from os/exec.
// As more functionality is needed, this can grow. Since command is a struct, we will have
// to replace fields with get/set method pairs.
type Interface interface {
	// Run runs the command to the completion.
	Run() error
	// CombinedOutput runs the command and returns its combined standard output
	// and standard error. This follows the pattern of package os/exec.
	CombinedOutput() ([]byte, error)
	// Output runs the command and returns standard output, but not standard err
	Output() ([]byte, error)
	SetDir(dir string)
	SetStdin(in io.Reader)
	SetStdout(out io.Writer)
	SetStderr(out io.Writer)
	SetEnv(env []string)

	// StdoutPipe and StderrPipe for getting the process' Stdout and Stderr as
	// Readers
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)

	// Start and Wait are for running a process non-blocking
	Start() error
	Wait() error

	// Stops the command by sending SIGTERM. It is not guaranteed the
	// process will stop before this function returns. If the process is not
	// responding, an internal timer function will send a SIGKILL to force
	// terminate after 10 seconds.
	Stop()
}

// Implements Interface in terms of really exec()ing.
type API struct{}

// Command is part of the Interface interface.
func (a *API) Command(cmd string, args ...string) Interface {
	return (*cmdWrapper)(osexec.Command(cmd, args...))
}

// CommandContext is part of the Interface interface.
func (a *API) CommandContext(ctx context.Context, cmd string, args ...string) Interface {
	return (*cmdWrapper)(osexec.CommandContext(ctx, cmd, args...))
}

// LookPath is part of the Interface interface
func (a *API) LookPath(file string) (string, error) {
	return osexec.LookPath(file)
}
