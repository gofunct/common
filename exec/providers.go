package exec

import (
	"github.com/gofunct/iio"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"os"
)

// New returns a new Interface which will os/exec to run commands.
func New(name string, i *iio.Service, args ...string) Interface {
	return &Executioner{
		Path:   name,
		Args:   args,
		Env:    viper.GetStringSlice("env"),
		Dir:    os.Getenv("PWD"),
		Stdin:  i.In(),
		Stdout: i.Out(),
		Stderr: i.Err(),
	}
}

var DefaultSet = wire.NewSet(
	New,
)
