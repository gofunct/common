package logging

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/svc/eventlog"

	"github.com/sirupsen/logrus"
)

func init() {
	setEventlogFormatter = func(l logger, name string, debugAsInfo bool) error {
		if name == "" {
			return fmt.Errorf("missing name parameter")
		}

		fmter, err := newEventlogger(name, debugAsInfo, l.entry.Logger.Formatter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating eventlog formatter: %v\n", err)
			l.Errorf("can't connect logger to eventlog: %v", err)
			return err
		}
		l.entry.Logger.Formatter = fmter
		return nil
	}
}

type eventlogger struct {
	log         *eventlog.Log
	debugAsInfo bool
	wrap        logrus.Formatter
}

func newEventlogger(name string, debugAsInfo bool, fmter logrus.Formatter) (*eventlogger, error) {
	logHandle, err := eventlog.Open(name)
	if err != nil {
		return nil, err
	}
	return &eventlogger{log: logHandle, debugAsInfo: debugAsInfo, wrap: fmter}, nil
}

func (s *eventlogger) Format(e *logrus.Entry) ([]byte, error) {
	data, err := s.wrap.Format(e)
	if err != nil {
		fmt.Fprintf(os.Stderr, "eventlogger: can't format entry: %v\n", err)
		return data, err
	}

	switch e.Level {
	case logrus.PanicLevel:
		fallthrough
	case logrus.FatalLevel:
		fallthrough
	case logrus.ErrorLevel:
		err = s.log.Error(102, e.Message)
	case logrus.WarnLevel:
		err = s.log.Warning(101, e.Message)
	case logrus.InfoLevel:
		err = s.log.Info(100, e.Message)
	case logrus.DebugLevel:
		if s.debugAsInfo {
			err = s.log.Info(100, e.Message)
		}
	default:
		err = s.log.Info(100, e.Message)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "eventlogger: can't send log to eventlog: %v\n", err)
	}

	return data, err
}
