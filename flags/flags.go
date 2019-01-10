package flags

import (
	"github.com/gofunct/common/config"
	"github.com/spf13/viper"
	"os"
)

type Interface interface {
	viper.FlagValue
}

type Cflag struct {
	changed bool
	key     string
	val     string
	bind    bool
}

func (f Cflag) SetName(k string) {
	f.key = k
}
func (f Cflag) SetVal(v string) {
	f.val = v
}

func (f Cflag) HasChanged() bool {
	return f.changed
}
func (f Cflag) IsEnv() bool {
	return f.bind
}

func (f Cflag) SetEnv() {
	if f.bind {
		os.Setenv(f.key, f.val)
	}
}

func (f Cflag) Name() string {
	return f.key
}

func (f Cflag) ValueString() string {

	return string(f.val)
}

func (f Cflag) ValueType() string { return "string" }

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

func NewFlagset() *Flagset {
	return &Flagset{}
}
