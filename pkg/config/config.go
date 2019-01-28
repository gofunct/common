package config

import (
	"github.com/gofunct/common/pkg/logger/zap"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	Viper = viper.New()
)

func init() {
	Viper.AutomaticEnv()
	Viper.AllowEmptyEnv(true)
	Viper.SetDefault("json_logs", false)
	Viper.SetDefault("loglevel", "debug")
	zap.LogE("Write config if it doesnt exist", Viper.SafeWriteConfig())
	Viper.AddConfigPath(os.Getenv("HOME" + filepath.Base(os.Getenv("PWD"))))
	Viper.SetConfigType("yaml")
	Viper.SetConfigName(".common")
	zap.LogE("Read in config", Viper.ReadInConfig())
	zap.LogE("Unmarshal Configuration variable from config file", Viper.Unmarshal(Configuration))
}

var Configuration *Config

// Config contains service configuration
type Config struct {
	Name        string `mapstructure:"pname"`
	Description string `mapstructure:"description"`
	Github      string `mapstructure:"github"`
	Project     string `mapstructure:"project"`
	Bin         string `mapstructure:"bin"`
	GitInit     bool   `mapstructure:"gitinit"`
	Contract    bool   `mapstructure:"contract"`
	GKE         struct {
		Enabled bool   `mapstructure:"gke.enabled"`
		Project string `mapstructure:"gke.project"`
		Zone    string `mapstructure:"gke.zone"`
		Cluster string `mapstructure:"gke.cluster"`
	}
	Storage struct {
		Enabled  bool `mapstructure:"storage.enabled"`
		Postgres bool `mapstructure:"storage.posgres"`
		MySQL    bool `mapstructure:"storage.mysql"`
		Config   struct {
			Driver      string `mapstructure:"storage.comfig.driver"`
			Host        string `mapstructure:"storage.comfig.driver"`
			Port        int    `mapstructure:"storage.comfig.driver"`
			Name        string `mapstructure:"storage.comfig.driver"`
			Username    string `mapstructure:"storage.comfig.driver"`
			Password    string `mapstructure:"storage.comfig.driver"`
			Connections struct {
				Max  int `mapstructure:"storage.comfig.connections.max"`
				Idle int `mapstructure:"storage.comfig.connections.idle"`
			}
		}
	}
	API struct {
		Enabled bool `mapstructure:"api.enabled"`
		GRPC    bool `mapstructure:"api.grpc"`
		Gateway bool `mapstructure:"api.gateway"`
		Config  struct {
			Port    int `mapstructure:"api.config.port"`
			Gateway struct {
				Port int `mapstructure:"api.config.port"`
			}
		}
	}
	Directories struct {
		Templates string `mapstructure:"directories.templates"`
		Service   string `mapstructure:"directories.service"`
	}
}
