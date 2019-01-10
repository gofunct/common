package flags

import (
	"github.com/spf13/viper"
	"os"
)

type Interface interface {
	viper.FlagValue
	SetEnv()
	IsEnv() bool
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
