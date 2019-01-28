package transport

import (
	"github.com/gofunct/common/pkg/transport/api"
	"github.com/gofunct/common/pkg/transport/engine"
)

func Serve(servers ...api.Server) error {
	s := engine.New(
		engine.WithDefaultLogger(),
		engine.WithServers(
			servers...,
		),
	)
	return s.Serve()
}
