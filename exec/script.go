package exec

import (
	"github.com/gofunct/iio"
	"io"
	osexec "os/exec"
)

// Wraps exec.Cmd so we can capture errors.
type Scripter struct {
	Name    string
	Usage   string
	IO      iio.Service
	RunFunc func(s *Scripter, args ...string) error
	OS      osexec.Cmd
}

func (s *Scripter) Output() []byte {
	return s.Output()
}

func (s *Scripter) Execute(args ...string) error {

	return s.RunFunc(s, args...)
}

func (s *Scripter) SetDir(dir string) {
	s.OS.Dir = dir
}

func (s *Scripter) SetStdin(in io.Reader) {
	s.IO.InR = in
}

func (s *Scripter) SetStdout(out io.Writer) {
	s.IO.OutW = out
}

func (s *Scripter) SetStderr(out io.Writer) {
	s.IO.ErrW = out
}

func (s *Scripter) SetEnv(env []string) {
	s.OS.Env = env
}

func (s *Scripter) GetArgs() []string {
	return s.OS.Args
}
