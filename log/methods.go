package log

import (
	"github.com/gofunct/common/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
)

func (l *Logger) Fatal(v ...interface{}) {
	l.Z.Fatal(utils.Red(v), zap.Any("fatal", v))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Z.Fatal(utils.Red(v), zap.Any("fatalf", v))
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.Z.Fatal(utils.Red(v), zap.Any("fatlln", v))
}

func (l *Logger) Flags() int {
	l.Z.Panic(utils.Yello("Flags are currently unsupported"))
	return 0
}

func (l *Logger) Output(calldepth int, s string) error {
	l.Z.Panic(utils.Yello("Flags are currently unsupported"))
	return nil
}

func (l *Logger) Panic(v ...interface{}) {
	l.Z.Panic(utils.Yello(v), zap.Any("panic", v))
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Z.Panic(utils.Yello(v), zap.Any("panicf", v))
}

func (l *Logger) Panicln(v ...interface{}) {
	l.Z.Panic(utils.Yello(v), zap.Any("panicln", v))
}

func (l *Logger) Prefix() string {
	return viper.GetString("log-prefix")
}

func (l *Logger) Print(v ...interface{}) {
	l.Z.Fatal(utils.Blue(v), zap.Any("print", v))
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Z.Fatal(utils.Blue(v), zap.Any("printf", v))
}

func (l *Logger) Println(v ...interface{}) {
	l.Z.Info(utils.Blue(v), zap.Any("println", v))
}

func (l *Logger) SetFlags(flag int) {
	l.Z.Panic(utils.Magenta("Setting flags is not supported"))
}

func (l *Logger) SetOutput(w io.Writer) {
	l.Z.Panic(utils.Magenta("changing output is currently unsupported"))
}

func (l *Logger) SetPrefix(prefix string) {
	viper.Set("log-prefix", prefix)
	l.Z.With(zap.String("prefix", prefix))
}
