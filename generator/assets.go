package generator

import (
	"github.com/jessevdk/go-assets"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Interface interface {
	Lister() map[string]*assets.File
	Get(string) string
	Add(path string) error
	Write(path string) error
	New(path string, data []byte) *assets.File
	Open(path string) (http.File, error)
}

//TODO: Add Interface
type API struct {
	Fs *assets.FileSystem
	G  *assets.Generator
}

func NewAssets(s *assets.FileSystem, g *assets.Generator) Interface {
	return API{
		Fs: s,
		G:  g,
	}
}

///////////////////////////////FS///////////////////////////////

func (a *API) Lister() map[string]*assets.File {
	return a.Fs.Files
}

func (a *API) Open(path string) (http.File, error) {
	f, err := a.Fs.Open(path)
	zap.L().Debug("Opening assetfs file", zap.String("path", path), zap.Error(err))
	return f, err
}

func (a *API) New(path string, data []byte) *assets.File {
	return a.Fs.NewFile(path, 0755, time.Now(), data)
}

///////////////////////////////Generator///////////////////////////////

func (a *API) Package() string {
	return a.G.PackageName
}

func (a *API) Prefix() string {
	return a.G.StripPrefix
}

func (a *API) Add(path string) error {
	err := a.G.Add(path)
	zap.L().Debug("Adding asset-generator file", zap.String("path", path), zap.Error(err))
	return err
}

func (a *API) Write(path string) error {
	err := a.G.Write(os.Stdout)
	zap.L().Debug("Writing asset-generator file", zap.String("path", path), zap.Error(err))
	return err
}

func (a *API) Get(f string) string {
	return a.G.VariableName
}

func (f *API) SortedRootAndTemplateFSFiles() []string {
	rootFiles := make([]string, 0, len(f.Fs.Files))
	tmplPaths := make([]string, 0, len(f.Fs.Files))
	for path, entry := range f.Fs.Files {
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
