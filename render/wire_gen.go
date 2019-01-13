// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package render

import (
	"github.com/google/wire"
)

// Injectors from providers.go:

func Inject() *Service {
	config := NewConfig()
	service := NewRenderer(config)
	return service
}

// providers.go:

var Provider = wire.NewSet(
	NewConfig,
	NewRenderer,
)
