package exec

import (
	"github.com/gofunct/iio"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io"
	"os/user"
)

// Wraps exec.Cmd so we can capture errors.
type Scripter struct {
	io *iio.Service
	v  *viper.Viper
}

func (s *Scripter) Output() []byte {
	return s.Output()
}

func (s *Scripter) SetStdin(in io.Reader) {
	s.io.InR = in
}

func (s *Scripter) SetStdout(out io.Writer) {
	s.io.OutW = out
}

func (s *Scripter) SetStderr(out io.Writer) {
	s.io.ErrW = out
}

func (s *Scripter) RequireRoot() error {
	u, err := user.Current()
	if err != nil {
		return errors.Wrapf(err, "%s\n", "failed to look up current user")
	}
	if u.Name != "root" {
		return errors.Wrapf(err, "%s\n", "root user is required")
	}

	return nil
}
