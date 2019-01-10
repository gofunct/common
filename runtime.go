package common

import (
	"github.com/gofunct/common/ask"
	"github.com/gofunct/common/config"
	"github.com/gofunct/common/exec"
	"github.com/gofunct/common/fs"
	"github.com/gofunct/common/log"
	"github.com/gofunct/iio"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"path/filepath"
	"strings"
)

type Runtime struct {
	exec.Service
	Object       interface{}
	Meta         map[string]interface{}
	Initializers []func() error
	Closers      []func() error
	WalkFunc     filepath.WalkFunc
	Handler      Handler
	Middleware   Middleware
	Runner       Runner
	Q            *ask.Service
	V            *config.Service
	IO           *iio.Service
	FS           *fs.Service
	L            *log.Service
}

func (r Runtime) SetObject(i interface{}) {
	r.Object = i
}

func (r Runtime) GetObject() interface{} {
	return r.Object
}

func (r *Runtime) UpdateMeta() error {
	if err := r.V.MergeConfigMap(r.V.AllSettings()); err != nil {
		return errors.Wrap(err, "failed to merge root and rervice config maps")
	}
	r.Meta = r.V.AllSettings()

	return nil
}

func (r *Runtime) Init() error {
	r.L.DebugC("initializing common runtime....\n")
	if err := r.UpdateMeta(); err != nil {
		return errors.WithStack(err)
	}
	if err := r.V.Unmarshal(r.Object); err != nil {
		return errors.WithStack(err)
	}

	for _, f := range r.Initializers {
		if err := f(); err != nil {
			return errors.WithStack(err)
		}
	}
	r.L.SuccessG("runtime has initialized successfully")
	return nil
}

func (r *Runtime) Generate(dir string) error {
	obj := r.GetObject()
	if err := r.V.Unmarshal(obj); err != nil {
		return err
	}
	for _, tmplPath := range r.FS.Generator.SortedRootAndTemplateFSFiles() {
		entry := r.FS.Generator.Fs.Files[tmplPath]
		path, err := TemplateString(strings.TrimSuffix(tmplPath, ".tmpl")).Compile(obj)
		if err != nil {
			return errors.Wrapf(err, "failed to parse path: %s", path)
		}
		absPath := filepath.Join(dir, path)
		dirPath := filepath.Dir(absPath)

		// create directory if not exists
		b, err := r.FS.Exists(dirPath)
		if err != nil {
			return errors.WithStack(err)
		}
		if !b {
			if err := r.FS.MakeDirAll(dirPath); err != nil {
				return errors.WithStack(err)
			}
		}

		// generate content
		body, err := TemplateString(string(entry.Data)).Compile(obj)
		if err != nil {
			return errors.Wrapf(err, "failed to generate %s", path)
		}

		err = afero.WriteFile(r.FS.Os, absPath, []byte(body), 0644)
		if err != nil {
			return errors.Wrapf(err, "failed to write %s", path)
		}

		r.L.DebugC(path[1:])
	}

	return nil
}
