package exec

import (
	"io"
	osexec "os/exec"
	"syscall"
	"time"
)

// Cmd is an interface that presents an Service that is very similar to Cmd from os/exec.
// As more functionality is needed, this can grow. Since command is s struct, we will have
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

	// Start and Wait are for running s process non-blocking
	Start() error
	Wait() error

	// Stops the command by sending SIGTERM. It is not guaranteed the
	// process will stop before this function returns. If the process is not
	// responding, an internal timer function will send s SIGKILL to force
	// terminate after 10 seconds.
	Stop()
}

// Wraps exec.Cmd so we can capture errors.
type Executioner osexec.Cmd

func (e *Executioner) SetDir(dir string) {
	e.Dir = dir
}

func (e *Executioner) SetStdin(in io.Reader) {
	e.Stdin = in
}

func (e *Executioner) SetStdout(out io.Writer) {
	e.Stdout = out
}

func (e *Executioner) SetStderr(out io.Writer) {
	e.Stderr = out
}

func (e *Executioner) SetEnv(env []string) {
	e.Env = env
}

func (e *Executioner) StdoutPipe() (io.ReadCloser, error) {
	r, err := (*osexec.Cmd)(e).StdoutPipe()
	return r, handleError(err)
}

func (e *Executioner) StderrPipe() (io.ReadCloser, error) {
	r, err := (*osexec.Cmd)(e).StderrPipe()
	return r, handleError(err)
}

func (e *Executioner) Start() error {
	err := (*osexec.Cmd)(e).Start()
	return handleError(err)
}

func (e *Executioner) Wait() error {
	err := (*osexec.Cmd)(e).Wait()
	return handleError(err)
}

// Run is part of the Cmd interface.
func (e *Executioner) Run() error {
	err := (*osexec.Cmd)(e).Run()
	return handleError(err)
}

// CombinedOutput is part of the Cmd interface.
func (e *Executioner) CombinedOutput() ([]byte, error) {
	out, err := (*osexec.Cmd)(e).CombinedOutput()
	return out, handleError(err)
}

func (e *Executioner) Output() ([]byte, error) {
	out, err := (*osexec.Cmd)(e).Output()
	return out, handleError(err)
}

// Stop is part of the Cmd interface.
func (e *Executioner) Stop() {
	c := (*osexec.Cmd)(e)

	if c.Process == nil {
		return
	}

	c.Process.Signal(syscall.SIGTERM)

	time.AfterFunc(10*time.Second, func() {
		if !c.ProcessState.Exited() {
			c.Process.Signal(syscall.SIGKILL)
		}
	})
}
