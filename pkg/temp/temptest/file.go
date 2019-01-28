package temptest

import (
	"bytes"
	"errors"
	"io"
)

// FakeFile is an implementation of a WriteCloser, that records what has
// been written in the file (in a bytes.Buffer) and if the file has been
// closed.
type FakeFile struct {
	Buffer bytes.Buffer
	Closed bool
}

var _ io.WriteCloser = &FakeFile{}

// Write appends the contents of p to the Buffer. If the file has
// already been closed, an error is returned.
func (f *FakeFile) Write(p []byte) (n int, err error) {
	if f.Closed {
		return 0, errors.New("can't write to closed FakeFile")
	}
	return f.Buffer.Write(p)
}

// Close records that the file has been closed. If the file has already
// been closed, an error is returned.
func (f *FakeFile) Close() error {
	if f.Closed {
		return errors.New("FakeFile was closed multiple times")
	}
	f.Closed = true
	return nil
}
