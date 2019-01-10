package log

import (
	"fmt"
	"github.com/gofunct/common/utils"
	"github.com/spf13/viper"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

type Logger struct {
	Z    *zap.Logger
	mode LoggingMode
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
func (l *Logger) AddLoggingFlags(cmd *cobra.Command) {
	var (
		debugEnabled, verboseEnabled bool
	)
	cmd.PersistentFlags().BoolVar(&debugEnabled, utils.Blue("debug"), false, utils.Blue("Debug level output"))
	cmd.PersistentFlags().BoolVarP(&verboseEnabled, utils.Blue("verbose"), "v", true, utils.Blue("Verbose loggingoutput"))

	cobra.OnInitialize(func() {
		switch {
		case debugEnabled:
			l.Z.With(
				zap.String("exec", cmd.Name()),
				zap.String("version", cmd.Version),
				zap.Bool("runnable", cmd.Runnable()))
			l.Debug()
		case verboseEnabled:
			l.Z.With(
				zap.String("exec", cmd.Name()),
				zap.String("version", cmd.Version),
				zap.Bool("runnable", cmd.Runnable()),
				zap.Any("meta", cmd.Annotations),
				zap.Bool("is-root", cmd.HasSubCommands()))
			l.VerboseLog()
		}
	})
}

// Debug sets a debug logger in global.
func (l *Logger) Debug() {
	logging = LoggingDebug
	l.ReplaceLoggerConfig(DebugLogConfig)
}

// Verbose sets a verbose logger in global.
func (l *Logger) VerboseLog() {
	logging = LoggingVerbose
	l.ReplaceLoggerConfig(VerboseLogConfig)
}

// IsDebug returns true if a debug logger is used.
func IsDebugLog() bool { return logging == LoggingDebug }

// IsVerbose returns true if a verbose logger is used.
func IsVerboseLog() bool { return logging == LoggingVerbose }

// Logging returns a current logging mode.
func (l *Logger) Mode() LoggingMode {
	return logging
}

func (l *Logger) ReplaceLoggerConfig(cfg zap.Config) {
	x, err := cfg.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize a debug logger: %v\n", err)
	}

	l.AddCloseFunc(func() { l.Z.Sync() })
	l.AddCloseFunc(zap.ReplaceGlobals(x))
}

func (l *Logger) Inititialize() {
	switch l.mode {
	case LoggingDebug:
		viper.Set("log-level", "debug")
		l.Z.With(
			zap.String("user", viper.GetString("USER")),
			zap.String("pwd", viper.GetString("PWD")))
		l.Debug()
	case LoggingVerbose:
		viper.Set("log-level", "verbpse")
		l.Z.With(
			zap.String("user", viper.GetString("USER")),
			zap.String("pwd", viper.GetString("PWD")),
			zap.String("version", runtime.Version()),
			zap.String("arch", runtime.GOARCH),
			zap.Int("cpus", runtime.NumCPU()),
			zap.Int("routines", runtime.NumGoroutine()))
		l.VerboseLog()
	}
}