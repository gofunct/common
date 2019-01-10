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
	ReadRemote() error
	ReadInConfig() error
	Unmarshal() error
	Set(string, interface{})
	GetString(string) string
	GetStringSlice(string) []string
	GetBool(string) bool
	GetObject() interface{}
	config.Provider
}

func NewConfig(vip *viper.Viper) Config {
	return config.API{
		V:        vip,
		Provider: &config.RemoteConfig{},
	}
}
