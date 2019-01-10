package common

import (
	"github.com/gofunct/common/flags"
)

type Flagger interface {
	flags.Interface
}

type FlagSet interface {
	flags.SetInterface
}

func NewFlag() Flagger {
	return &flags.Cflag{}
}

func NewFlagSet(name string) FlagSet {
	return &flags.Flagset{
		Name: name,
	}
}
