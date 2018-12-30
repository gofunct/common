package svcgen

import (
	"github.com/gofunct/common/proto/protoc"
	"github.com/gofunct/common/proto/service/params"
)

type CreateAppFunc func(*gencmd.Command) (*App, error)

type App struct {
	ProtocWrapper protoc.Wrapper
	ParamsBuilder params.Builder
}
