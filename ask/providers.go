//+build wireinject

package ask

import (
	"github.com/gofunct/iio"
	"github.com/google/wire"
	"github.com/tcnksm/go-input"
)

func NewDefault() *Service {
	return &Service{
		Q: input.DefaultUI(),
	}
}

func New(i *iio.Service) *Service {
	return &Service{
		Q: &input.UI{
			Writer: i.Out(),
			Reader: i.In(),
		},
	}
}

var Providers = wire.NewSet(
	iio.Inject,
	New,
)

func Inject() (*Service, error) {
	wire.Build(Providers)
	return &Service{}, nil
}
