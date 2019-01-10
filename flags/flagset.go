package flags

import (
	"github.com/gofunct/common/config"
	"github.com/spf13/viper"
)

type SetInterface interface {
	viper.FlagValueSet
	AddFlag(v ...Interface)
	SetPrefix(p string)
	FlagsFromConfig(c config.API)
}

type Flagset struct {
	Name   string
	flags  []Interface
	prefix string
}

func (f *Flagset) VisitAll(fn func(viper.FlagValue)) {
	for _, k := range f.flags {
		fn(k)
	}
}

func (f *Flagset) SetPrefix(p string) {
	f.prefix = p
}

func (f *Flagset) AddFlag(v ...Interface) {
	for _, val := range v {
		f.flags = append(f.flags, val)
	}
}

func (f *Flagset) FlagsFromConfig(c config.API) {
	for k, v := range c.AllSettings() {
		for _, x := range f.flags {
			if k != x.Name() {
				f.flags = append(f.flags, &Cflag{
					changed: false,
					key:     k,
					val:     v.(string),
					bind:    true,
				})
			}
		}
	}
}

func (f *Flagset) AddCflags(c ...Cflag) {
	for _, v := range c {
		f.flags = append(f.flags, v)
	}
}
