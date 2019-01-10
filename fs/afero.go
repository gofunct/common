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

type API struct {
	*afero.Afero
}

func (a *API) Create(name string) (afero.File, error) {

	return a.Fs.Create(name)

}

///////////////////////////CHECK///////////////////////////

func (a *API) CheckFilepathHasPrefix(path string, prefix string) bool {
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
func (r *API) CheckIfCmdDir(name string) bool {
	name = filepath.Base(name)
	for _, cmdDir := range []string{"cmd", "cmds", "command", "commands"} {
		if name == cmdDir {
			return true
		}
	}
	return false
}

func (a *API) CheckIfThisIsDir(path string) (bool, error) {
	return a.IsDir(path)
}

func (a *API) CheckIfFileContainThis(filename string, this []byte) (bool, error) {
	return a.FileContainsBytes(filename, this)
}

func (a *API) CheckIfThisDirEmpty(path string) bool {
	b, err := a.IsEmpty(path)
	zap.L().Fatal("Checking if directory is empty", zap.String("path", path), zap.Error(err))
	return b
}

///////////////////////////MAKE///////////////////////////

func (a *API) MakeDir(path string) error {
	err := a.MkdirAll(path, 0755)
	zap.L().Debug("Making All Directories", zap.String("path", path), zap.Error(err))
	return errors.Wrapf(err, "failed to create %q directory", path)
}

func (a *API) MakeTempFile(dir, prefix string) (afero.File, error) {
	f, err := a.TempFile(dir, prefix)
	zap.L().Debug("Making Temporary File", zap.String("dir", dir), zap.String("prefix", prefix), zap.Error(err))
	return f, err
}

func (a *API) MakeTempDir(dir, prefix string) (string, error) {
	s, err := a.TempDir(dir, prefix)
	zap.L().Debug("Making Temporary Directory", zap.String("dir", dir), zap.String("prefix", prefix), zap.Error(err))
	return s, err
}

///////////////////////////WRITE///////////////////////////

func (a *API) WriteToFile(filename string, data []byte) error {
	err := a.WriteFile(filename, data, 0755)
	zap.L().Debug("Writing to File", zap.String("filename", filename), zap.ByteString("data", data), zap.Error(err))
	return err
}

func (a *API) WriteToReader(path string, r io.Reader) error {
	err := a.WriteReader(path, r)
	zap.L().Debug("Writing to File", zap.String("path", path), zap.Any("reader", r), zap.Error(err))
	return err
}

///////////////////////////READ///////////////////////////

func (a *API) ReadFromDir(path string) ([]os.FileInfo, error) {
	i, err := a.ReadDir(path)
	zap.L().Debug("Reading directory", zap.String("path", path), zap.Error(err))

	return i, err
}

func (a *API) ReadFromFile(path string) ([]byte, error) {
	b, err := a.ReadFile(path)
	zap.L().Debug("Reading file", zap.String("path", path), zap.Error(err))
	return b, err
}

func (a *API) OpenFile(path string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := a.Open(path)
	zap.L().Debug("Opening file", zap.String("path", path), zap.Error(err))
	return f, err
}

///////////////////////////WALK///////////////////////////

func (a *API) WalkPath(path string, walkFn filepath.WalkFunc) error {
	err := a.Walk(path, walkFn)
	zap.L().Debug("Walking path with func", zap.String("path", path), zap.Error(err))
	return err
}

///////////////////////////LIST///////////////////////////

///////////////////////////DELETE///////////////////////////

func (a *API) Remove(path string) error {
	err := a.Remove(path)
	zap.L().Debug("Removing file", zap.String("path", path), zap.Error(err))
	return err
}

///////////////////////////OTHER///////////////////////////

func (a *API) Rename(old, new string) error {
	err := a.Rename(old, new)
	zap.L().Debug("Renaming", zap.String("old", old), zap.String("new", new), zap.Error(err))
	return err
}

func (a *API) ChangePermissions(path string, o os.FileMode) error {
	err := a.Chmod(path, o)
	zap.L().Debug("Changing permissions", zap.String("path", path), zap.Any("file-mode", o), zap.Error(err))
	return err
}

func (a *API) Stat(name string) (os.FileInfo, error) {
	o, err := a.Stat(name)
	zap.L().Debug("Changing permissions", zap.String("name", name), zap.Error(err))
	return o, err
}

// exists checks if a file or directory exists.
func (f *API) Exists(path string) (bool, error) {
	if path == "" {
		return false, nil
	}
	_, err := f.Stat(path)
	if err == nil {
		return true, nil
	}
	if !os.IsNotExist(err) {
		return true, err
	}
	return false, err
}

// findCmdDir checks if base of absPath is cmd dir and returns it or
// looks for existing cmd dir in absPath.
func (f *API) FindCmdDir(absPath string) string {
	absExists, err := f.Exists(absPath)
	if err != nil {
		zap.L().Fatal("failed to check if abs path exists", zap.Error(err))
	}

	if !absExists || f.CheckIfThisDirEmpty(absPath) {
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
func (f *API) FindPackage(packageName string) string {
	if packageName == "" {
		return ""
	}

	for _, srcPath := range srcPaths {
		packagePath := filepath.Join(srcPath, packageName)
		b, err := f.Exists(packagePath)
		if err != nil {
			zap.L().Fatal("failed to check if package path exists", zap.Error(err))
		}
		if b {
			return packagePath
		}
	}

	return ""
}

// trimSrcPath trims at the beginning of absPath the srcPath.
func (f *API) TrimSrcPath(absPath, srcPath string) string {
	relPath, err := filepath.Rel(srcPath, absPath)
	if err != nil {
		zap.L().Fatal("failed to get abs from src path", zap.Error(err))
	}
	return relPath
}

// isCmdDir checks if base of name is one of cmdDir.
func (f *API) IsCmdDir(name string) bool {
	name = filepath.Base(name)
	for _, cmdDir := range []string{"cmd", "cmds", "command", "commands"} {
		if name == cmdDir {
			return true
		}
	}
	return false
}

func (f *API) FilepathHasPrefix(path string, prefix string) bool {
	if len(path) <= len(prefix) {
		return false
	}
	if runtime.GOOS == "windows" {
		// Paths in windows are case-insensitive.
		return strings.EqualFold(path[0:len(prefix)], prefix)
	}
	return path[0:len(prefix)] == prefix

}

func (f *API) GetCurrentUser() string {
	usr, _ := user.Current()
	return usr.Username
}

func (f *API) NewBasePathFs(dir string) {
	f.Fs = afero.NewBasePathFs(afero.NewOsFs(), dir)
}

func (f *API) NewMemMapOs() {
	f.Fs = afero.NewMemMapFs()
}
