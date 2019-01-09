package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type Config struct {
	v    *viper.Viper
	Meta map[string]interface{}
}

func (c *Config) SetFs(a afero.Fs) {
	c.v.SetFs(a)

}

func (c *Config) Unmarshal(i interface{}) error {
	if err := c.v.Unmarshal(i); err != nil {
		return errors.Wrapf(err, "failed to unmarshal object")
	}

	return nil
}

func (c *Config) Bytes(i interface{}) ([]byte, error) {
	b, ok := i.([]byte)
	if ok {
		return b, nil
	}
	return nil, errors.New("type does not contain any bytes")

}

func (c *Config) GetEnv() ([]string, error) {
	s := c.v.GetStringSlice("ENV")
	if len(s) == 0 {
		return nil, errors.New("configurator was queried for ENV, but no values were returned")
	}
	return s, nil
}

func (c *Config) Value(key string) interface{} {
	return c.v.Get(key)
}

func (c *Config) MergeMeta(v map[string]interface{}) {
	c.v.MergeConfigMap(v)
}

func (c *Config) GetMeta() map[string]interface{} {
	return c.v.AllSettings()
}

func (c *Config) Debug() {
	viper.Debug()
}
