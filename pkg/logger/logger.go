package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
}

func WithWriter(writer io.Writer) *logrus.Logger {
	Logger.Out = writer
	return Logger
}

func WithLevel() {
	switch viper.GetString("loglevel") {
	case "debug":
		Logger.Level = logrus.DebugLevel
	case "warning":
		Logger.Level = logrus.WarnLevel
	case "info":
		Logger.Level = logrus.InfoLevel
	default:
		Logger.Level = logrus.DebugLevel
	}
}

// LogF logs fatal "msg: err" in case of error
func LogF(msg string, err error) {
	if err != nil {
		Logger.Fatalf("%s: %s", msg, err)
	}
}

// LogE logs error "msg: err" in case of error
func LogE(msg string, err error) {
	if err != nil {
		Logger.Printf("%s: %s", msg, err)
	}
}

// LogE logs error "msg: err" in case of error
func Context(k, v string) *logrus.Entry {
	return Logger.WithField(k, v)
}
