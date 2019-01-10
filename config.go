package common

import (
	"github.com/gofunct/common/config"
	"github.com/spf13/viper"
)

type Config interface {
	AllSettings() map[string]interface{}
	Debug()
	MergeConfigMap(map[string]interface{}) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
	SetObject(i interface{})
}

func NewEncoder() Config {
	return config.API{
		Viper: viper.New(),
	}
}
