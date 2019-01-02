//+build wireinject

package svcgen

import (
	"github.com/gofunct/gogen/pkg/cli"
	"github.com/gofunct/gogen/pkg/gencmd"
	"github.com/gofunct/gogen/pkg/protoc"
	"github.com/google/wire"
)

func NewApp(*gencmd.Command) (*App, error) {
	wire.Build(
		Set,
		gencmd.Set,
		cli.UIInstance,
		protoc.WrapperSet,
	)
	return nil, nil
}
