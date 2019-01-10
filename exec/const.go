package exec

import (
	osexec "os/exec"
)

var (
	ErrExecutableNotFound           = osexec.ErrNotFound
	_                     ExitError = CodeExitError{}
	_                     Interface = &Executioner{}

// ErrExecutableNotFound is returned if the executable is not found.
)
