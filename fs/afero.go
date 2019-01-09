package fs

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	AferoOS  = afero.NewOsFs()
	AferoMem = afero.NewMemMapFs()
	srcPaths []string
)

type Interface interface {
	DirExists(path string) (bool, error)
	Exists(path string) (bool, error)
	FileContainsBytes(filename string, subslice []byte) (bool, error)
	GetTempDir(subPath string) string
	IsDir(path string) (bool, error)
	IsEmpty(path string) (bool, error)
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	SafeWriteReader(path string, r io.Reader) (err error)
	TempDir(dir, prefix string) (name string, err error)
	TempFile(dir, prefix string) (f afero.File, err error)
	Walk(root string, walkFn filepath.WalkFunc) error
	WriteFile(filename string, data []byte, perm os.FileMode) error
	WriteReader(path string, r io.Reader) (err error)
}

type Api struct {
	*afero.Afero
}

func (a *Api) Create(name string) (afero.File, error) {

	return a.Fs.Create(name)

}

///////////////////////////CHECK///////////////////////////

func (c *Api) CheckFilepathHasPrefix(path string, prefix string) bool {
	if len(path) <= len(prefix) {
		return false
	}
	if runtime.GOOS == "windows" {
		// Paths in windows are case-insensitive.
		return strings.EqualFold(path[0:len(prefix)], prefix)
	}
	return path[0:len(prefix)] == prefix

}

// isCmdDir checks if base of name is one of cmdDir.
func (r *Api) CheckIfCmdDir(name string) bool {
	name = filepath.Base(name)
	for _, cmdDir := range []string{"cmd", "cmds", "command", "commands"} {
		if name == cmdDir {
			return true
		}
	}
	return false
}

func (a *Api) CheckIfThisIsDir(path string) (bool, error) {
	return a.IsDir(path)
}

func (a *Api) CheckIfFileContainThis(filename string, this []byte) (bool, error) {
	return a.FileContainsBytes(filename, this)
}

func (a *Api) CheckIfThisDirEmpty(path string) bool {
	b, err := a.IsEmpty(path)
	zap.L().Fatal("Checking if directory is empty", zap.String("path", path), zap.Error(err))
	return b
}

func (a *Api) FindAllThisPattern(pattern string) ([]string, error) {
	f, err := afero.Glob(a, pattern)
	zap.L().Debug("Finding all files with pattern", zap.String("pattern", pattern), zap.Error(err))
	return f, err
}

func (a *Api) FindAllProtoFiles() ([]string, error) {
	f, err := afero.Glob(a, "*.proto")
	zap.L().Debug("Finding all proto files", zap.Error(err))
	return f, err
}

func (a *Api) FindAllGoFiles() ([]string, error) {
	f, err := afero.Glob(a, "*.go")
	zap.L().Debug("Finding all go files", zap.Error(err))
	return f, err
}

func (a *Api) FindAllYamlFiles() ([]string, error) {
	f, err := afero.Glob(a, "*.yaml")
	zap.L().Debug("Finding all yaml files", zap.Error(err))
	return f, err
}

func (a *Api) FindAllJsonFiles() ([]string, error) {
	f, err := afero.Glob(a, "*.json")
	zap.L().Debug("Finding all json files", zap.Error(err))
	return f, err
}

func (a *Api) FindAllMdFiles() ([]string, error) {
	f, err := afero.Glob(a, "*.md")
	zap.L().Debug("Finding all markdown files", zap.Error(err))
	return f, err
}

func (a *Api) FindAllPBFiles() ([]string, error) {
	f, err := afero.Glob(a, "*.pb.go")
	zap.L().Debug("Finding all generated protobuf files", zap.Error(err))
	return f, err
}

func (a *Api) FindAllShelllFiles() ([]string, error) {
	f, err := afero.Glob(a, "*.sh")
	zap.L().Debug("Finding all shell files", zap.Error(err))
	return f, err
}

///////////////////////////MAKE///////////////////////////

func (a *Api) MakeDir(path string) error {
	err := a.MkdirAll(path, 0755)
	zap.L().Debug("Making All Directories", zap.String("path", path), zap.Error(err))
	return errors.Wrapf(err, "failed to create %q directory", path)
}

func (a *Api) MakeTempFile(dir, prefix string) (afero.File, error) {
	f, err := a.TempFile(dir, prefix)
	zap.L().Debug("Making Temporary File", zap.String("dir", dir), zap.String("prefix", prefix), zap.Error(err))
	return f, err
}

func (a *Api) MakeTempDir(dir, prefix string) (string, error) {
	s, err := a.TempDir(dir, prefix)
	zap.L().Debug("Making Temporary Directory", zap.String("dir", dir), zap.String("prefix", prefix), zap.Error(err))
	return s, err
}

///////////////////////////WRITE///////////////////////////

func (a *Api) WriteToFile(filename string, data []byte) error {
	err := a.WriteFile(filename, data, 0755)
	zap.L().Debug("Writing to File", zap.String("filename", filename), zap.ByteString("data", data), zap.Error(err))
	return err
}

func (a *Api) WriteToReader(path string, r io.Reader) error {
	err := a.WriteReader(path, r)
	zap.L().Debug("Writing to File", zap.String("path", path), zap.Any("reader", r), zap.Error(err))
	return err
}

///////////////////////////READ///////////////////////////

func (a *Api) ReadFromDir(path string) ([]os.FileInfo, error) {
	i, err := a.ReadDir(path)
	zap.L().Debug("Reading directory", zap.String("path", path), zap.Error(err))

	return i, err
}

func (a *Api) ReadFromFile(path string) ([]byte, error) {
	b, err := a.ReadFile(path)
	zap.L().Debug("Reading file", zap.String("path", path), zap.Error(err))
	return b, err
}

func (a *Api) OpenFile(path string) (afero.File, error) {
	f, err := a.Open(path)
	zap.L().Debug("Opening file", zap.String("path", path), zap.Error(err))
	return f, err
}

///////////////////////////WALK///////////////////////////

func (a *Api) WalkPath(path string, walkFn filepath.WalkFunc) error {
	err := a.Walk(path, walkFn)
	zap.L().Debug("Walking path with func", zap.String("path", path), zap.Error(err))
	return err
}

///////////////////////////LIST///////////////////////////

///////////////////////////DELETE///////////////////////////

func (a *Api) Remove(path string) error {
	err := a.Remove(path)
	zap.L().Debug("Removing file", zap.String("path", path), zap.Error(err))
	return err
}

///////////////////////////OTHER///////////////////////////

func (a *Api) Rename(old, new string) error {
	err := a.Rename(old, new)
	zap.L().Debug("Renaming", zap.String("old", old), zap.String("new", new), zap.Error(err))
	return err
}

func (a *Api) ChangePermissions(path string, o os.FileMode) error {
	err := a.Chmod(path, o)
	zap.L().Debug("Changing permissions", zap.String("path", path), zap.Any("file-mode", o), zap.Error(err))
	return err
}

func (a *Api) Stat(name string) (os.FileInfo, error) {
	o, err := a.Stat(name)
	zap.L().Debug("Changing permissions", zap.String("name", name), zap.Error(err))
	return o, err
}

// exists checks if a file or directory exists.
func (f *Api) Exists(path string) bool {
	if path == "" {
		return false
	}
	_, err := f.Stat(path)
	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		zap.L().Fatal("file or directory already exists", zap.Error(err))
	}
	return false
}

// findCmdDir checks if base of absPath is cmd dir and returns it or
// looks for existing cmd dir in absPath.
func (f *Api) FindCmdDir(absPath string) string {
	if !f.Exists(absPath) || f.CheckIfThisDirEmpty(absPath) {
		return "cmd"
	}

	if f.CheckIfCmdDir(absPath) {
		return filepath.Base(absPath)
	}

	files, _ := filepath.Glob(filepath.Join(absPath, "c*"))
	for _, file := range files {
		if f.CheckIfCmdDir(file) {
			return filepath.Base(file)
		}
	}

	return "cmd"
}

// findPackage returns full path to existing go package in GOPATHs.
func (f *Api) FindPackage(packageName string) string {
	if packageName == "" {
		return ""
	}

	for _, srcPath := range srcPaths {
		packagePath := filepath.Join(srcPath, packageName)
		if f.Exists(packagePath) {
			return packagePath
		}
	}

	return ""
}

// trimSrcPath trims at the beginning of absPath the srcPath.
func (f *Api) TrimSrcPath(absPath, srcPath string) string {
	relPath, err := filepath.Rel(srcPath, absPath)
	if err != nil {
		zap.L().Fatal("failed to get abs from src path", zap.Error(err))
	}
	return relPath
}

// isCmdDir checks if base of name is one of cmdDir.
func (f *Api) IsCmdDir(name string) bool {
	name = filepath.Base(name)
	for _, cmdDir := range []string{"cmd", "cmds", "command", "commands"} {
		if name == cmdDir {
			return true
		}
	}
	return false
}

func (f *Api) FilepathHasPrefix(path string, prefix string) bool {
	if len(path) <= len(prefix) {
		return false
	}
	if runtime.GOOS == "windows" {
		// Paths in windows are case-insensitive.
		return strings.EqualFold(path[0:len(prefix)], prefix)
	}
	return path[0:len(prefix)] == prefix

}

func (f *Api) GetCurrentUser() string {
	usr, _ := user.Current()
	return usr.Username
}
