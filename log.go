package common

import (
	"fmt"
	"github.com/gofunct/common/log"
	"go.uber.org/zap"
	"io"
	"os"
)

//Log is an interface that is compatible with the stdLib log.Logger
type Log interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Flags() int
	Output(calldepth int, s string) error
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Prefix() string
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	SetFlags(flag int)
	SetOutput(w io.Writer)
	SetPrefix(prefix string)
}

func NewZapProductionLog() Log {
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

func NewZapDevelopmentLog() Log {
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
