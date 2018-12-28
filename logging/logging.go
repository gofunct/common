
package logging

import (
	"github.com/dixonwille/wlog"
	kitlog "github.com/go-kit/kit/log"
	"log"
	"os"
	"strconv"
)

type Logger struct{
	UI 		wlog.UI
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
			UI:            wlog.New(os.Stdin, os.Stdout, os.Stderr),
		},
		KitLog: logger,
	}

	log.SetOutput(kitlog.NewStdlibAdapter(l.KitLog))
	return l
}


func (l *Logger) FatalUI(args ...interface{}) {
	for _, arg := range args {
		switch x := arg.(type) {
		case string:
			l.UI.Log(x)
		case error:
			l.UI.Log(x.Error())
		case int:
			l.UI.Log(string(x))
		case bool:
			l.UI.Log(strconv.FormatBool(x))
		case byte:
			l.UI.Log(string(x))
		}
	}
}

func (l *Logger) InfoUI(args ...interface{}) {
	for _, arg := range args {
		switch x := arg.(type) {
		case string:
			l.UI.Info(x)
		case error:
			l.UI.Info(x.Error())
		case int:
			l.UI.Info(string(x))
		case bool:
			l.UI.Info(strconv.FormatBool(x))
		case byte:
			l.UI.Info(string(x))
		}
	}
}

func (l *Logger) ErrorUI(args ...interface{}) {
	for _, arg := range args {
		switch x := arg.(type) {
		case string:
			l.UI.Error(x)
		case error:
			l.UI.Error(x.Error())
		case int:
			l.UI.Error(string(x))
		case bool:
			l.UI.Error(strconv.FormatBool(x))
		case byte:
			l.UI.Error(string(x))
		}
	}
}

func (l *Logger) OutputUI(args ...interface{}) {
	for _, arg := range args {
		switch x := arg.(type) {
		case string:
			l.UI.Output(x)
		case error:
			l.UI.Output(x.Error())
		case int:
			l.UI.Output(string(x))
		case bool:
			l.UI.Output(strconv.FormatBool(x))
		case byte:
			l.UI.Output(string(x))
		}
	}
}

func (l *Logger) RunningUI(args ...interface{}) {
	for _, arg := range args {
		switch x := arg.(type) {
		case string:
			l.UI.Running(x)
		case error:
			l.UI.Running(x.Error())
		case int:
			l.UI.Running(string(x))
		case bool:
			l.UI.Running(strconv.FormatBool(x))
		case byte:
			l.UI.Running(string(x))
		}
	}
}


func (l *Logger) WarnUI(args ...interface{}) {
	for _, arg := range args {
		switch x := arg.(type) {
		case string:
			l.UI.Warn(x)
		case error:
			l.UI.Warn(x.Error())
		case int:
			l.UI.Warn(string(x))
		case bool:
			l.UI.Warn(strconv.FormatBool(x))
		case byte:
			l.UI.Warn(string(x))
		}
	}
}

func (l *Logger) SuccessUI(args ...interface{}) {
	for _, arg := range args {
		switch x := arg.(type) {
		case string:
			l.UI.Success(x)
		case error:
			l.UI.Success(x.Error())
		case int:
			l.UI.Success(string(x))
		case bool:
			l.UI.Success(strconv.FormatBool(x))
		case byte:
			l.UI.Success(string(x))
		}
	}
}


func (l *Logger) DebugUI(args ...interface{}) {
	for _, arg := range args {
		switch x := arg.(type) {
		case string:
			l.UI.Info(x)
		case error:
			l.UI.Info(x.Error())
		case int:
			l.UI.Info(string(x))
		case bool:
			l.UI.Info(strconv.FormatBool(x))
		case byte:
			l.UI.Info(string(x))
		}
	}
}
