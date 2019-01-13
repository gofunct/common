//+build wireinject

package fs

import (
	"github.com/google/wire"
	"github.com/spf13/afero"
	"os"
)

func New() *Service {
	var rt = RootDir(os.Getenv("PWD"))
	osFs := afero.NewOsFs()
	baseFs := afero.NewBasePathFs(osFs, rt.String())
	s := &Service{
		Os: &afero.Afero{
			Fs: baseFs,
		},
		Root: rt,
	}
	s.HttpFs = afero.NewHttpFs(s.Os)

	return s
}

func NewFromHOME() *Service {
	var rt = RootDir(os.Getenv("HOME"))
	osFs := afero.NewOsFs()
	baseFs := afero.NewBasePathFs(osFs, rt.String())
	s := &Service{
		Os: &afero.Afero{
			Fs: baseFs,
		},
		Root: rt,
	}
	s.HttpFs = afero.NewHttpFs(s.Os)

	return s
}

func NewFromGOPATH() *Service {
	var rt = RootDir(os.Getenv("GOPATH"))
	osFs := afero.NewOsFs()
	baseFs := afero.NewBasePathFs(osFs, rt.String())
	s := &Service{
		Os: &afero.Afero{
			Fs: baseFs,
		},
		Root: rt,
	}
	s.HttpFs = afero.NewHttpFs(s.Os)

	return s
}

var Provider = wire.NewSet(
	New,
)

func Inject() (*Service, error) {
	wire.Build(Provider)
	return &Service{}, nil
}
