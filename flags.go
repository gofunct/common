package common

import (
	"github.com/gofunct/common/flags"
	"github.com/spf13/viper"
)

type Flagger interface {
	viper.FlagValueSet
}

func NewFlagSet() viper.FlagValueSet {
	//TODO: add initializer for every type in their source directory
	return &flags.Flagset{}
}

func NewFlag() viper.FlagValue {
	return &flags.Cflag{}
}
