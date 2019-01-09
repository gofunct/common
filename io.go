package common

import (
	iio "github.com/gofunct/common/io"
	"github.com/mattn/go-colorable"
	"io"
	"os"
)

// IO contains an input reader, an output writer and an error writer.
type IO interface {
	In() io.Reader
	Out() io.Writer
	Err() io.Writer
}

// Stdio returns a standard IO object.
func NewStdIO() IO {
	io := &iio.IOContainer{
		InR:  os.Stdin,
		OutW: colorable.NewColorableStdout(),
		ErrW: colorable.NewColorableStderr(),
	}

	return io
}
