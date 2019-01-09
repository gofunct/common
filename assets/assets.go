package assets

import (
	"github.com/jessevdk/go-assets"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

//TODO: Add Interface
type Assets struct {
	fs *assets.FileSystem
	g  *assets.Generator
}

func NewAssets() *Assets {
	return &Assets{
		fs: &assets.FileSystem{},
		g:  &assets.Generator{},
	}
}

///////////////////////////////FS///////////////////////////////

func (a *Assets) ListFSFiles() map[string]*assets.File {
	return a.fs.Files
}

func (a *Assets) OpenFSFile(path string) (http.File, error) {
	f, err := a.fs.Open(path)
	zap.L().Debug("Opening assetfs file", zap.String("path", path), zap.Error(err))
	return f, err
}

func (a *Assets) NewFSFile(path string, data []byte) *assets.File {
	return a.fs.NewFile(path, 0755, time.Now(), data)
}

///////////////////////////////Generator///////////////////////////////

func (a *Assets) GetFSPackage() string {
	return a.g.PackageName
}

func (a *Assets) GetFSPrefix() string {
	return a.g.StripPrefix
}

func (a *Assets) AddFileToFS(path string) error {
	err := a.g.Add(path)
	zap.L().Debug("Adding asset-generator file", zap.String("path", path), zap.Error(err))
	return err
}

func (a *Assets) WriteFSFile(path string) error {
	err := a.g.Write(os.Stdout)
	zap.L().Debug("Writing asset-generator file", zap.String("path", path), zap.Error(err))
	return err
}

func (a *Assets) GetFSVariable(f string) string {
	return a.g.VariableName
}

func (f *Assets) SortedRootAndTemplateFSFiles() []string {
	rootFiles := make([]string, 0, len(f.fs.Files))
	tmplPaths := make([]string, 0, len(f.fs.Files))
	for path, entry := range f.fs.Files {
		if entry.IsDir() {
			continue
		}
		if strings.Count(entry.Path[1:], "/") == 0 {
			rootFiles = append(rootFiles, path)
		} else {
			tmplPaths = append(tmplPaths, path)
		}
	}
	sort.Strings(rootFiles)
	sort.Strings(tmplPaths)
	return append(rootFiles, tmplPaths...)
}
