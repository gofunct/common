package temptest

import (
	"errors"
	"fmt"
	"io"

	"k8s.io/utils/temp"
)

// FakeDir implements a Directory that is not backed on the
// filesystem. This is useful for testing since the created "files" are
// simple bytes.Buffer that can be inspected.
type FakeDir struct {
	Files   map[string]*FakeFile
	Deleted bool
}

var _ temp.Directory = &FakeDir{}

// NewFile returns a new FakeFile if the filename doesn't exist already.
// This function will fail if the directory has already been deleted.
func (d *FakeDir) NewFile(name string) (io.WriteCloser, error) {
	if d.Deleted {
		return nil, errors.New("can't create file in deleted FakeDir")
	}
	if d.Files == nil {
		d.Files = map[string]*FakeFile{}
	}
	f := d.Files[name]
	if f != nil {
		return nil, fmt.Errorf(
			"FakeDir already has file named %q",
			name,
		)
	}
	f = &FakeFile{}
	d.Files[name] = f
	return f, nil
}

// Delete doesn't remove anything, but records that the directory has
// been deleted.
func (d *FakeDir) Delete() error {
	if d.Deleted {
		return errors.New("failed to re-delete FakeDir")
	}
	d.Deleted = true
	return nil
}
