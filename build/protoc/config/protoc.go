package config

import (
	"github.com/gofunct/common/build/protoc/plugin"
	"github.com/gofunct/common/fs"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"
)

// Service stores setting params related protoc.
type Service struct {
	ImportDirs []string `mapstructure:"import_dirs"`
	ProtosDir  string   `mapstructure:"protos_dir"`
	OutDir     string   `mapstructure:"out_dir"`
	Plugins    []*plugin.Plugin
}

// ProtoFiles returns .proto file paths.
func (c *Service) ProtoFiles(f fs.Service) ([]string, error) {
	paths := []string{}
	protosDir := filepath.Join(f.Root.String(), c.ProtosDir)
	err := afero.Walk(f.Os, protosDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}
		if !info.IsDir() && filepath.Ext(path) == ".proto" {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, errors.WithStack(err)
}

// OutDirOf returns a directory path of protoc result output path for given proto file.
func (c *Service) OutDirOf(rootDir string, protoPath string) (string, error) {
	protosDir := filepath.Join(rootDir, c.ProtosDir)
	relProtoDir, err := filepath.Rel(protosDir, filepath.Dir(protoPath))
	if strings.Contains(relProtoDir, "..") {
		return "", errors.Errorf(".proto files should be included in %s", c.ProtosDir)
	}
	if err != nil {
		return "", errors.Wrapf(err, ".proto files should be included in %s", c.ProtosDir)
	}

	return filepath.Join(c.OutDir, relProtoDir), nil
}

// Commands returns protoc command and arguments for given proto file.
func (c *Service) Commands(rootDir, protoPath string) ([][]string, error) {
	cmds := make([][]string, 0, len(c.Plugins))
	relProtoPath, _ := filepath.Rel(rootDir, protoPath)

	outDir, err := c.OutDirOf(rootDir, protoPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, p := range c.Plugins {
		args := []string{}
		args = append(args, "-I", filepath.Dir(relProtoPath))
		for _, dir := range c.ImportDirs {
			args = append(args, "-I", dir)
		}
		args = append(args, p.ToProtocArg(outDir))
		args = append(args, relProtoPath)
		cmds = append(cmds, append([]string{"protoc"}, args...))
	}

	return cmds, nil
}
