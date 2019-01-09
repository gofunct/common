package exec

import (
	osexec "os/exec"
	"syscall"
)

// ExitError is an interface that presents an API similar to os.ProcessState, which is
// what ExitError from os/exec is. This is designed to make testing a bit easier and
// probably loses some of the cross-platform properties of the underlying library.
type ExitError interface {
	String() string
	Error() string
	Exited() bool
	ExitStatus() int
}

// ExitErrorWrapper is an implementation of ExitError in terms of os/exec ExitError.
// Note: standard exec.ExitError is type *os.ProcessState, which already implements Exited().
type ExitErrorWrapper struct {
	*osexec.ExitError
}

// CodeExitError is an implementation of ExitError consisting of an error object
// and an exit code (the upper bits of os.exec.ExitStatus).
type CodeExitError struct {
	Err  error
	Code int
}

// ExitStatus is part of the ExitError interface.
func (eew ExitErrorWrapper) ExitStatus() int {
	ws, ok := eew.Sys().(syscall.WaitStatus)
	if !ok {
		panic("can't call ExitStatus() on a non-WaitStatus exitErrorWrapper")
	}
	return ws.ExitStatus()
}

func (e CodeExitError) Error() string {
	return e.Err.Error()
}

func (e CodeExitError) String() string {
	return e.Err.Error()
}

// Exited is to check if the process has finished
func (e CodeExitError) Exited() bool {
	return true
}

// ExitStatus is for checking the error code
func (e CodeExitError) ExitStatus() int {
	return e.Code
}

func handleError(err error) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case *osexec.ExitError:
		return &ExitErrorWrapper{e}
	case *osexec.Error:
		if e.Err == osexec.ErrNotFound {
			return ErrExecutableNotFound
		}
	}

	return err
}
