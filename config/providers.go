package config

import (
	"github.com/google/wire"
)

func New() (*Service, error) {
	s := &Service{}

	if err := s.Init(); err != nil {
		return nil, err
	}

	return s, nil
}

var DefaultSet = wire.NewSet(
	New,
)
