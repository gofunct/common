package io

import (
	"bytes"
	"io"
)

// IOContainer is a basic implementation of the IO interface.
type IOContainer struct {
	InR  io.Reader
	OutW io.Writer
	ErrW io.Writer
}

func (i *IOContainer) In() io.Reader  { return i.InR }
func (i *IOContainer) Out() io.Writer { return i.OutW }
func (i *IOContainer) Err() io.Writer { return i.ErrW }

// FakeIO is a fake implementation of the IO interface using `bytes.Buffer`s.
type FakeIO struct {
	InBuf  *bytes.Buffer
	OutBuf *bytes.Buffer
	ErrBuf *bytes.Buffer
}

func (i *FakeIO) In() io.Reader  { return i.InBuf }
func (i *FakeIO) Out() io.Writer { return i.OutBuf }
func (i *FakeIO) Err() io.Writer { return i.ErrBuf }
