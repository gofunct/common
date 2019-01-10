package common

import (
	"github.com/gofunct/common/config"
	"github.com/spf13/viper"
	"github.com/google/wire"
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
	viper.RemoteProvider
    SetPath(string)
    SetProvider(string)
    SetEndpoint(string)
    SetSecretKeyring(string)
}

func NewConfig(vip *viper.Viper) Config {
	return config.API{
		V:        vip,
	}
}

var ConfigSet = wire.NewSet(
	NewConfig,
)


