//+build wireinject

package common

import (
	"github.com/gofunct/common/ask"
	"github.com/gofunct/common/config"
	"github.com/gofunct/common/exec"
	"github.com/gofunct/common/fs"
	"github.com/gofunct/common/log"
	"github.com/gofunct/common/render"
	"github.com/gofunct/iio"
	"github.com/google/wire"
	"github.com/izumin5210/grapi/pkg/protoc"
	"github.com/jessevdk/go-assets"
	"github.com/spf13/cobra"
	"path/filepath"
)

var DefaultSet = wire.NewSet(
	ask.DefaultSet,
	config.DefaultSet,
	fs.DefaultSet,
	iio.Set,
	exec.DefaultSet,
	log.VerboseSet,
	protoc.WrapperSet,
	render.Set,
)

func NewAsk() (*ask.Service, error) {
	wire.Build(DefaultSet)
	return &ask.Service{}, nil
}
func NewConfig() (*config.Service, error) {
	wire.Build(DefaultSet)
	return &config.Service{}, nil
}
func NewFs(walkFunc filepath.WalkFunc, f *assets.FileSystem) (*fs.Service, error) {
	wire.Build(DefaultSet)
	return &fs.Service{}, nil
}

func NewIO() (*iio.Service, error) {
	wire.Build(DefaultSet)
	return &iio.Service{}, nil
}

func NewLog() (*log.Service, error) {
	wire.Build(DefaultSet)
	return &log.Service{}, nil
}

func NewCommander(s *exec.Scripter) *cobra.Command {
	wire.Build(DefaultSet)
	return exec.NewCommander(s)
}

////////

func NewVerboseLog() (*log.Service, error) {
	wire.Build(DefaultSet)
	return &log.Service{}, nil
}

func NewRenderer() *render.Renderer {
	wire.Build(DefaultSet)
	return &render.Renderer{}
}
