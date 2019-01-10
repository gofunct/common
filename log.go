package common

import (
	"fmt"
	"github.com/gofunct/common/log"
	"go.uber.org/zap"
	"os"
)

type Logger interface {
	log.Log
}

func NewZapProductionLog() Logger {
	l, err := zap.NewProduction()
	if err != nil {
		fmt.Println("failed to setup logger")
		os.Exit(1)
	}
	zap.ReplaceGlobals(l)
	x := &log.Logger{
		Z: l,
	}
	x.Debug()
	x.Inititialize()
	return x
}

func NewZapDevelopmentLog() Logger {
	l, err := zap.NewProduction()
	if err != nil {
		fmt.Println("failed to setup logger")
		os.Exit(1)
	}
	zap.ReplaceGlobals(l)
	x := &log.Logger{
		Z: l,
	}
	x.VerboseLog()
	x.Inititialize()
	return x
}
