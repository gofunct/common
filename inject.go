package common

import (
    "github.com/spf13/viper"
	"github.com/google/wire"
)

func New(v *viper.Viper) Config {
	wire.Build(ConfigSet)
	return nil
}