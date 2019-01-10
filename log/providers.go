package log

import (
	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func NewVerbose() (*Service, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	zap.ReplaceGlobals(l)
	return &Service{
		Z:    l,
		mode: LoggingVerbose,
	}, nil
}

func NewDebug() (*Service, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	zap.ReplaceGlobals(l)
	return &Service{
		Z:    l,
		mode: LoggingDebug,
	}, nil
}

var DefaultSet = wire.NewSet(
	NewDebug,
)

var VerboseSet = wire.NewSet(
	NewVerbose,
)
