package common

import (
	"bytes"
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
	io := &ioContainer{
		InR:  os.Stdin,
		OutW: colorable.NewColorableStdout(),
		ErrW: colorable.NewColorableStderr(),
	}

	return io
}

// IOContainer is a basic implementation of the IO interface.
type ioContainer struct {
	InR  io.Reader
	OutW io.Writer
	ErrW io.Writer
}

func (i *ioContainer) In() io.Reader  { return i.InR }
func (i *ioContainer) Out() io.Writer { return i.OutW }
func (i *ioContainer) Err() io.Writer { return i.ErrW }

// FakeIO is a fake implementation of the IO interface using `bytes.Buffer`s.
type FakeIO struct {
	InBuf  *bytes.Buffer
	OutBuf *bytes.Buffer
	ErrBuf *bytes.Buffer
}

func (i *FakeIO) In() io.Reader  { return i.InBuf }
func (i *FakeIO) Out() io.Writer { return i.OutBuf }
func (i *FakeIO) Err() io.Writer { return i.ErrBuf }
