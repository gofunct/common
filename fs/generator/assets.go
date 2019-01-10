package generator

import (
	"github.com/jessevdk/go-assets"
	"go.uber.org/zap"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Service struct {
	Fs *assets.FileSystem
}

///////////////////////////////FS///////////////////////////////

func (s *Service) Lister() map[string]*assets.File {
	return s.Fs.Files
}

func (s *Service) OpenFile(path string) (http.File, error) {
	f, err := s.Fs.Open(path)
	zap.L().Debug("Opening assetfs file", zap.String("path", path), zap.Error(err))
	return f, err
}

func (s *Service) NewFile(path string, data []byte) *assets.File {
	return s.Fs.NewFile(path, 0755, time.Now(), data)
}

func (f *Service) SortedRootAndTemplateFSFiles() []string {
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
