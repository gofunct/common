package hack

import (
	"github.com/gofunct/iio"
	"github.com/pkg/errors"
	"io"
	"os/user"
)

// Wraps exec.Cmd so we can capture errors.
type Service struct {
	io *iio.Service
}

func (s *Service) Output() []byte {
	return s.Output()
}

func (s *Service) SetStdin(in io.Reader) {
	s.io.InR = in
}

func (s *Service) SetStdout(out io.Writer) {
	s.io.OutW = out
}

func (s *Service) SetStderr(out io.Writer) {
	s.io.ErrW = out
}

func (s *Service) RequireRoot() error {
	u, err := user.Current()
	if err != nil {
		return errors.Wrapf(err, "%s\n", "failed to look up current user")
	}
	if u.Name != "root" {
		return errors.Wrapf(err, "%s\n", "root user is required")
	}

	return nil
}
