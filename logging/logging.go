
package logging

import (
	"github.com/dixonwille/wlog"
	kitlog "github.com/go-kit/kit/log"
	"log"
	"os"
)


type Logger struct{
	wlog.UI
	Logger kitlog.Logger
}

func NewLogger() *Logger {
	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))

	return &Logger{
		UI: &wlog.PrefixUI{
			LogPrefix:     ":speech_balloon:",
			OutputPrefix:  ":boom:",
			SuccessPrefix: ":white_check_mark:",
			InfoPrefix:    ":wave:",
			ErrorPrefix:   ":x:",
			WarnPrefix:    ":warning:",
			RunningPrefix: ":zap:",
			AskPrefix:     ":interrobang:",
			UI:            wlog.New(os.Stdin, kitlog.NewStdlibAdapter(logger), kitlog.NewStdlibAdapter(logger)),
		},
		Logger: logger,
	}
}

func (l *Logger) UseStdLog() {
	log.SetOutput(kitlog.NewStdlibAdapter(l.Logger))
	log.Print("Logger initialized")
}






