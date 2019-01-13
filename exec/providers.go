package exec

import (
	"github.com/gofunct/iio"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// New returns a new Interface which will os/exec to run commands.
func NewScripter() *Scripter {
	return &Scripter{
		io: iio.NewStdIO(),
		v:  viper.GetViper(),
	}
}

var DefaultSet = wire.NewSet(
	NewScripter,
)
