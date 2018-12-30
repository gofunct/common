package logging

import (
	"fmt"
	"github.com/dixonwille/wlog"
	"github.com/gofunct/common/io"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var ui = wlog.New(os.Stdin, os.Stdout, os.Stderr)

var prefix = &wlog.PrefixUI{
	LogPrefix:     emoji.Sprint(":speech_balloon:"),
	OutputPrefix:   emoji.Sprint(":boom:"),
	SuccessPrefix:  emoji.Sprint(":white_check_mark:"),
	InfoPrefix:     emoji.Sprint(":wave:"),
	ErrorPrefix:    emoji.Sprint(":x:"),
	WarnPrefix:     emoji.Sprint(":grimacing:"),
	RunningPrefix:  emoji.Sprint(":fire:"),
	AskPrefix:      emoji.Sprint(":question:"),
	UI:            ui,
}


type Messenger struct{
	UI 		wlog.UI
}

func NewMessenger() *Messenger {
	m := &Messenger{
		UI: prefix,
	}
	m.AddColor()

	return m
}


func (m *Messenger) AddColor() {
		m.UI = wlog.AddColor(wlog.Green, wlog.Red, wlog.BrightBlue, wlog.Blue, wlog.Yellow, wlog.BrightMagenta, wlog.Yellow, wlog.BrightGreen, wlog.BrightRed, l.UI)
}


// LoggingMode represents a logging configuration specification.
type LoggingMode int

// LoggingMode values
const (
	LoggingNop LoggingMode = iota
	LoggingVerbose
	LoggingDebug
)

var (
	logging = LoggingNop

	// DebugLogConfig is used to generate a *zap.Logger for debug mode.
	DebugLogConfig = func() zap.Config {
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		cfg.DisableStacktrace = true
		return cfg
	}()
	// VerboseLogConfig is used to generate a *zap.Logger for verbose mode.
	VerboseLogConfig = func() zap.Config {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Local().Format("2006-01-02 15:04:05 MST"))
		}
		return cfg
	}()
)

// AddLoggingFlags sets "--debug" and "--verbose" flags to the given *cobra.Command instance.
func AddLoggingFlags(cmd *cobra.Command) {
	var (
		debugEnabled, verboseEnabled bool
	)

	cmd.PersistentFlags().BoolVar(
		&debugEnabled,
		"debug",
		false,
		fmt.Sprintf("Debug level output"),
	)
	cmd.PersistentFlags().BoolVarP(
		&verboseEnabled,
		"verbose",
		"v",
		false,
		fmt.Sprintf("Verbose level output"),
	)

	cobra.OnInitialize(func() {
		switch {
		case debugEnabled:
			Debug()
		case verboseEnabled:
			Verbose()
		}
	})
}

// Debug sets a debug logger in global.
func Debug() {
	logging = LoggingDebug
	replaceLogger(DebugLogConfig)
}

// Verbose sets a verbose logger in global.
func Verbose() {
	logging = LoggingVerbose
	replaceLogger(VerboseLogConfig)
}

// IsDebug returns true if a debug logger is used.
func IsDebug() bool { return logging == LoggingDebug }

// IsVerbose returns true if a verbose logger is used.
func IsVerbose() bool { return logging == LoggingVerbose }

// Logging returns a current logging mode.
func Logging() LoggingMode { return logging }

func replaceLogger(cfg zap.Config) {
	l, err := cfg.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize a debug logger: %v\n", err)
	}

	io.AddCloseFunc(func() { l.Sync() })
	io.AddCloseFunc(zap.ReplaceGlobals(l))
}
