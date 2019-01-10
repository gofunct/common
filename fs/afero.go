package fs

import (
	"github.com/gofunct/common/fs/generator"
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

type Service struct {
	Os        *afero.Afero
	HttpFs    *afero.HttpFs
	Root      RootDir
	Generator *generator.Service
	WalkFunc  filepath.WalkFunc
}

///////////////////////////CHECK///////////////////////////

func (s *Service) CheckFilepathHasPrefix(path string, prefix string) bool {
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
func (r *Service) CheckIfCmdDir(name string) bool {
	name = filepath.Base(name)
	for _, cmdDir := range []string{"cmd", "cmds", "command", "commands"} {
		if name == cmdDir {
			return true
		}
	}
	return false
}

///////////////////////////MAKE///////////////////////////

func (s *Service) MakeDirAll(path string) error {
	err := s.Os.MkdirAll(path, 0755)
	zap.L().Debug("Making All Directories", zap.String("path", path), zap.Error(err))
	return errors.Wrapf(err, "failed to create %q directory", path)
}

func (s *Service) MakeTempFile(dir, prefix string) (afero.File, error) {
	f, err := s.Os.TempFile(dir, prefix)
	zap.L().Debug("Making Temporary File", zap.String("dir", dir), zap.String("prefix", prefix), zap.Error(err))
	return f, err
}

func (s *Service) MakeTempDir(dir, prefix string) (string, error) {
	out, err := s.Os.TempDir(dir, prefix)
	zap.L().Debug("Making Temporary Directory", zap.String("dir", dir), zap.String("prefix", prefix), zap.Error(err))
	return out, err
}

///////////////////////////WRITE///////////////////////////

func (s *Service) WriteToFile(filename string, data []byte) error {
	err := s.Os.WriteFile(filename, data, 0755)
	zap.L().Debug("Writing to File", zap.String("filename", filename), zap.ByteString("data", data), zap.Error(err))
	return err
}

func (s *Service) WriteToReader(path string, r io.Reader) error {
	err := s.Os.WriteReader(path, r)
	zap.L().Debug("Writing to File", zap.String("path", path), zap.Any("reader", r), zap.Error(err))
	return err
}

///////////////////////////READ///////////////////////////

func (s *Service) ReadFromDir(path string) ([]os.FileInfo, error) {
	i, err := s.Os.ReadDir(path)
	zap.L().Debug("Reading directory", zap.String("path", path), zap.Error(err))

	return i, err
}

func (s *Service) ReadFromFile(path string) ([]byte, error) {
	b, err := s.Os.ReadFile(path)
	zap.L().Debug("Reading file", zap.String("path", path), zap.Error(err))
	return b, err
}

func (s *Service) OpenFile(path string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := s.Os.Open(path)
	zap.L().Debug("Opening file", zap.String("path", path), zap.Error(err))
	return f, err
}

///////////////////////////WALK///////////////////////////

func (s *Service) WalkRoot() error {
	b, err := s.Os.IsEmpty(s.Root.String())
	if err != nil {
		return errors.Wrap(err, "failed to check if root.root directory is empty")
	}
	if b {
		return errors.New("cannot walk an empty root.root directory")
	}
	b, err = s.Os.Exists(s.Root.TemplatesDir())
	if err != nil {
		return errors.Wrap(err, "failed to check if root.rppt directory exists")
	}
	if !b {
		return errors.New("root.root directory is nonexistent")
	}
	err = s.Os.Walk(s.Root.String(), s.WalkFunc)
	zap.L().Debug("Walking path with func", zap.String("path", s.Root.String()), zap.Error(err))
	return err
}

func (s *Service) WalkTemplates() error {
	b, err := s.Os.IsEmpty(s.Root.TemplatesDir())
	if err != nil {
		return errors.Wrap(err, "failed to check if root.templates directory is empty")
	}
	if b {
		return errors.New("cannot walk an empty root.templates directory")
	}
	b, err = s.Os.Exists(s.Root.TemplatesDir())
	if err != nil {
		return errors.Wrap(err, "failed to check if root.templates directory exists")
	}
	if !b {
		return errors.New("root.templates directory is nonexistent")
	}

	err = s.Os.Walk(s.Root.TemplatesDir(), s.WalkFunc)
	zap.L().Debug("Walking path with func", zap.String("path", s.Root.TemplatesDir()), zap.Error(err))
	return err
}
func (s *Service) WalkBin() error {
	b, err := s.Os.IsEmpty(s.Root.BinDir())
	if err != nil {
		return errors.Wrap(err, "failed to check if root.bin directory is empty")
	}
	if b {
		return errors.New("cannot walk an empty root.bin directory")
	}
	b, err = s.Os.Exists(s.Root.BinDir())
	if err != nil {
		return errors.Wrap(err, "failed to check if root.bin directory exists")
	}
	if !b {
		return errors.New("root.bin directory is nonexistent")
	}

	err = s.Os.Walk(s.Root.BinDir(), s.WalkFunc)
	zap.L().Debug("Walking path with func", zap.String("path", s.Root.BinDir()), zap.Error(err))
	return err
}

func (s *Service) WalkCmd() error {
	b, err := s.Os.IsEmpty(s.Root.CmdDir())
	if err != nil {
		return errors.Wrap(err, "failed to check if root.cmd directory is empty")
	}
	if b {
		return errors.New("cannot walk an empty root.cmd directory")
	}
	b, err = s.Os.Exists(s.Root.CmdDir())
	if err != nil {
		return errors.Wrap(err, "failed to check if root.cmd directory exists")
	}
	if !b {
		return errors.New("root.cmd directory is nonexistent")
	}

	err = s.Os.Walk(s.Root.CmdDir(), s.WalkFunc)
	zap.L().Debug("Walking path with func", zap.String("path", s.Root.CmdDir()), zap.Error(err))
	return err
}

func (s *Service) WalkProto() error {
	b, err := s.Os.IsEmpty(s.Root.ProtoDir())
	if err != nil {
		return errors.Wrap(err, "failed to check if root.proto directory is empty")
	}
	if b {
		return errors.New("cannot walk an empty root.proto directory")
	}
	b, err = s.Os.Exists(s.Root.ProtoDir())
	if err != nil {
		return errors.Wrap(err, "failed to check if root.proto directory exists")
	}
	if !b {
		return errors.New("root.proto directory is nonexistent")
	}
	err = s.Os.Walk(s.Root.ProtoDir(), s.WalkFunc)
	zap.L().Debug("Walking path with func", zap.String("path", s.Root.ProtoDir()), zap.Error(err))
	return err
}

func (s *Service) Walk(p string) error {
	b, err := s.Os.IsEmpty(p)
	if err != nil {
		return errors.Wrap(err, "failed to check if directory is empty")
	}
	if b {
		return errors.New("cannot walk an empty directory")
	}
	b, err = s.Os.Exists(p)
	if err != nil {
		return errors.Wrap(err, "failed to check if directory exists")
	}
	if !b {
		return errors.New("cannot walk a non-existent directory")
	}
	err = s.Os.Walk(p, s.WalkFunc)
	zap.L().Debug("Walking path with func", zap.String("path", p), zap.Error(err))
	return err
}

///////////////////////////LIST///////////////////////////

///////////////////////////DELETE///////////////////////////

func (s *Service) Remove(path string) error {
	err := s.Remove(path)
	zap.L().Debug("Removing file", zap.String("path", path), zap.Error(err))
	return err
}

///////////////////////////OTHER///////////////////////////

func (s *Service) Rename(old, new string) error {
	err := s.Rename(old, new)
	zap.L().Debug("Renaming", zap.String("old", old), zap.String("new", new), zap.Error(err))
	return err
}

func (s *Service) ChangePermissions(path string, o os.FileMode) error {
	err := s.Os.Chmod(path, o)
	zap.L().Debug("Changing permissions", zap.String("path", path), zap.Any("file-mode", o), zap.Error(err))
	return err
}

func (s *Service) Stat(name string) (os.FileInfo, error) {
	o, err := s.Stat(name)
	zap.L().Debug("Changing permissions", zap.String("name", name), zap.Error(err))
	return o, err
}

// exists checks if s file or directory exists.
func (f *Service) Exists(path string) (bool, error) {
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

// findPackage returns full path to existing go package in GOPATHs.
func (f *Service) FindPackage(packageName string) string {
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
func (f *Service) TrimSrcPath(absPath, srcPath string) string {
	relPath, err := filepath.Rel(srcPath, absPath)
	if err != nil {
		zap.L().Fatal("failed to get abs from src path", zap.Error(err))
	}
	return relPath
}

// isCmdDir checks if base of name is one of cmdDir.
func (f *Service) IsCmdDir(name string) bool {
	name = filepath.Base(name)
	for _, cmdDir := range []string{"cmd", "cmds", "command", "commands"} {
		if name == cmdDir {
			return true
		}
	}
	return false
}

func (f *Service) FilepathHasPrefix(path string, prefix string) bool {
	if len(path) <= len(prefix) {
		return false
	}
	if runtime.GOOS == "windows" {
		// Paths in windows are case-insensitive.
		return strings.EqualFold(path[0:len(prefix)], prefix)
	}
	return path[0:len(prefix)] == prefix

}

func (f *Service) GetCurrentUser() string {
	usr, _ := user.Current()
	return usr.Username
}
