package generator

import (
	"github.com/google/wire"
	"github.com/jessevdk/go-assets"
)

func New(fs *assets.FileSystem) *Service {
	return &Service{
		Fs: fs,
	}
}

var DefaultSet = wire.NewSet(
	New,
)
