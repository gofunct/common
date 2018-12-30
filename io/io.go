package io

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

// IOContainer is a basic implementation of the IO interface.
type IOContainer struct {
	InR  io.Reader
	OutW io.Writer
	ErrW io.Writer
}

func (i *IOContainer) In() io.Reader  { return i.InR }
func (i *IOContainer) Out() io.Writer { return i.OutW }
func (i *IOContainer) Err() io.Writer { return i.ErrW }

// Stdio returns a standard IO object.
func Stdio() IO {
	io := &IOContainer{
		InR:  os.Stdin,
		OutW: colorable.NewColorableStdout(),
		ErrW: colorable.NewColorableStderr(),
	}

	return io
}

// FakeIO is a fake implementation of the IO interface using `bytes.Buffer`s.
type FakeIO struct {
	InBuf  *bytes.Buffer
	OutBuf *bytes.Buffer
	ErrBuf *bytes.Buffer
}

func (i *FakeIO) In() io.Reader  { return i.InBuf }
func (i *FakeIO) Out() io.Writer { return i.OutBuf }
func (i *FakeIO) Err() io.Writer { return i.ErrBuf }

// NewFakeIO returns a new FakeIO object.
func NewFakeIO() *FakeIO {
	return &FakeIO{
		InBuf:  new(bytes.Buffer),
		OutBuf: new(bytes.Buffer),
		ErrBuf: new(bytes.Buffer),
	}
}
