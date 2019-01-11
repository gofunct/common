package exec

import (
	"context"
	"io"
	osexec "os/exec"
)

// Wraps exec.Cmd so we can capture errors.
type Scripter struct {
	osexec.Cmd
	Script ScriptFunc
}

func (e *Scripter) SetDir(dir string) {
	e.Dir = dir
}

func (e *Scripter) SetStdin(in io.Reader) {
	e.Stdin = in
}

func (e *Scripter) SetStdout(out io.Writer) {
	e.Stdout = out
}

func (e *Scripter) SetStderr(out io.Writer) {
	e.Stderr = out
}

func (e *Scripter) SetEnv(env []string) {
	e.Env = env
}

func (e *Scripter) Execute() error {
	return e.Run()
}

func tester() ScriptFunc {
	return func(cmd osexec.Cmd) Script {
		return func(ctx context.Context, name string, args ...string) error {
			return nil
		}
	}
}
