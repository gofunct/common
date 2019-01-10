package fs

import (
	"github.com/gofunct/common/fs/generator"
	"github.com/google/wire"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
)

func New(f filepath.WalkFunc) *Service {
	var rt = RootDir(os.Getenv("PWD"))
	osFs := afero.NewOsFs()
	baseFs := afero.NewBasePathFs(osFs, rt.String())
	s := &Service{
		Os: &afero.Afero{
			Fs: baseFs,
		},
		WalkFunc: f,
		Root:     rt,
	}
	s.HttpFs = afero.NewHttpFs(s.Os)

	return s
}

func NewFromHOME(f filepath.WalkFunc) *Service {
	var rt = RootDir(os.Getenv("HOME"))
	osFs := afero.NewOsFs()
	baseFs := afero.NewBasePathFs(osFs, rt.String())
	s := &Service{
		Os: &afero.Afero{
			Fs: baseFs,
		},
		WalkFunc: f,
		Root:     rt,
	}
	s.HttpFs = afero.NewHttpFs(s.Os)

	return s
}

func NewFromGOPATH(f filepath.WalkFunc) *Service {
	var rt = RootDir(os.Getenv("GOPATH"))
	osFs := afero.NewOsFs()
	baseFs := afero.NewBasePathFs(osFs, rt.String())
	s := &Service{
		Os: &afero.Afero{
			Fs: baseFs,
		},
		WalkFunc: f,
		Root:     rt,
	}
	s.HttpFs = afero.NewHttpFs(s.Os)

	return s
}

func NewWithGenerator(f filepath.WalkFunc, g *generator.Service) *Service {
	var rt = RootDir(os.Getenv("PWD"))
	osFs := afero.NewOsFs()
	baseFs := afero.NewBasePathFs(osFs, rt.String())
	s := &Service{
		Os: &afero.Afero{
			Fs: baseFs,
		},
		WalkFunc: f,
		Root:     rt,
	}
	s.HttpFs = afero.NewHttpFs(s.Os)
	s.Generator = generator.New(g.Fs)

	return s
}

var DefaultSet = wire.NewSet(
	generator.DefaultSet,
	NewWithGenerator,
)

var HomeSet = wire.NewSet(
	generator.DefaultSet,
	NewFromHOME,
)

var GopathSet = wire.NewSet(
	generator.DefaultSet,
	NewFromGOPATH,
)
