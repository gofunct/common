package logging

import (
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
}

func (l *Logger) Debug(args ...interface{}) {
	l.DebugUI(args)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.DebugUI(args)
}

func (l *Logger) Debugf(f string, args ...interface{}) {
	l.DebugUI(args)
}

func (l *Logger) Info(args ...interface{}) {
	l.InfoUI(args)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.InfoUI(args)
}

func (l *Logger) Infof(f string, args ...interface{}) {
	l.InfoUI(args)
}

func (l *Logger) Warn(args ...interface{}) {
	l.WarnUI(args)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.WarnUI(args)
}

func (l *Logger) Warnf(r string, args ...interface{}) {
	l.WarnUI(args)
}

func (l *Logger) Error(args ...interface{}) {
	l.ErrorUI(args)
	os.Exit(1)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.ErrorUI(args)
	os.Exit(1)
}

func (l *Logger) Errorf(f string, args ...interface{}) {
	l.ErrorUI(args)
	os.Exit(1)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.FatalUI(args)
	os.Exit(1)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.FatalUI(args)
	os.Exit(1)
}

func (l *Logger) Fatalf(f string, args ...interface{}) {
	l.FatalUI(args)
	os.Exit(1)
}

func (l *Logger) With(key string, value interface{}) *Logger {
	l.KitLog.Log(key, value)
	return l
}

