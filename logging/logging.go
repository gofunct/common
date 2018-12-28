
package logging

import (
	"github.com/dixonwille/wlog"
	kitlog "github.com/go-kit/kit/log"
	"log"
	"os"
)

type UILogger interface {
	Debug(...interface{})
	Debugln(...interface{})
	Debugf(string, ...interface{})

	Info(...interface{})
	Infoln(...interface{})
	Infof(string, ...interface{})

	Warn(...interface{})
	Warnln(...interface{})
	Warnf(string, ...interface{})

	Error(...interface{})
	Errorln(...interface{})
	Errorf(string, ...interface{})

	Fatal(...interface{})
	Fatalln(...interface{})
	Fatalf(string, ...interface{})

	With(key string, value interface{}) Logger

	SetFormat(string) error
	SetLevel(string) error
}

type Logger struct{
	wlog.UI
	KitLog kitlog.Logger
}

func NewLogger() *Logger {
	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))

	l := &Logger{
		UI: &wlog.PrefixUI{
			LogPrefix:     ":speech_balloon:",
			OutputPrefix:  ":boom:",
			SuccessPrefix: ":white_check_mark:",
			InfoPrefix:    ":wave:",
			ErrorPrefix:   ":x:",
			WarnPrefix:    ":warning:",
			RunningPrefix: ":zap:",
			AskPrefix:     ":interrobang:",
			UI:            wlog.New(os.Stdin, kitlog.NewStdlibAdapter(logger), os.Stderr),
		},
		KitLog: logger,
	}
	log.SetOutput(kitlog.NewStdlibAdapter(l.KitLog))
	log.Print("Logger initialized")

	return l
}
