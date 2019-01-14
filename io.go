package common

import (
	"io"
)

// IO is a basic implementation of the IO interface.
type IO struct {
	InR  io.Reader
	OutW io.Writer
	ErrW io.Writer
}

func (i *IO) In() io.Reader  { return i.InR }
func (i *IO) Out() io.Writer { return i.OutW }
func (i *IO) Err() io.Writer { return i.ErrW }
