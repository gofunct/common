//+build wireinject

package render

import "github.com/google/wire"

var Provider = wire.NewSet(
	NewConfig,
	NewRenderer,
)

func Inject() *Service {
	wire.Build(Provider)
	return &Service{}
}
