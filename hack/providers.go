//+build wireinject

package hack

import (
	"github.com/gofunct/iio"
	"github.com/google/wire"
)

// New returns a new Interface which will os/exec to run commands.
func New() *Service {
	return &Service{
		io: iio.NewStdIO(),
	}
}

var Provider = wire.NewSet(
	New,
)

func Inject() *Service {
	wire.Build(Provider)
	return &Service{}
}
