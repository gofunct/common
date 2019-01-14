package common

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
	"os/user"
	"runtime"
	"github.com/mitchellh/go-homedir"
)
var (
	// Find home directory.
	home, _ = homedir.Dir()
)
type Config struct {
	Bucket         string `mapstructure:"bucket"`
	DbHost         string `mapstructure:"db_host"`
	DbName         string `mapstructure:"db_name"`
	DbUser         string `mapstructure:"db_user"`
	DbPassword     string `mapstructure:"db_password"`
	CloudSqlRegion string `mapstructure:"cloud_sql_region"`
	Deploy         string `mapstructure:"deploy"`
	Lis            string `mapstructure:"lis"`
	Source         string `mapstructure:"source"`
	Container      string `mapstructure:"container"`
	SchemaPath     string `mapstructure:"schema_path"`
	RolesPath      string `mapstructure:"roles_path"`
	GenSource      string `mapstructure:"gen_source"`
	GenDest        string `mapstructure:"gen_dest"`
	GenPkgName     string `mapstructure:"gen_pkg_name"`
	*viper.Viper
}

func NewConfig() (*Config, error) {
	s := &Config{}
	s.Viper = viper.New()

	ex, err := os.Executable()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	s.AddConfigPath(home)
	s.SetConfigName(".common")
	s.AutomaticEnv()

	s.SetDefault("home", home)
	s.SetDefault("executable", ex)
	gr, err := os.Getgroups()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	s.SetDefault("groups", gr)
	host, err := os.Hostname()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.SetDefault("env", os.Environ())
	s.SetDefault("uid", os.Getuid())
	s.SetDefault("args", os.Args)
	s.SetDefault("host_name", host)
	s.SetDefault("pid", os.Getpid())
	s.SetDefault("goarch", runtime.GOARCH)
	s.SetDefault("compiler", runtime.Compiler)
	s.SetDefault("runtime_version", runtime.Version())
	s.SetDefault("goos", runtime.GOOS)
	s.SetDefault("goroot", runtime.GOROOT())
	usr, _ := user.Current()
	s.SetDefault("user", usr)

	if err := s.Unmarshal(s); err != nil {
		return nil, errors.WithStack(err)
	}

	// If a config file is found, read it in.
	if err := s.ReadInConfig(); err != nil {
		log.Println("failed to read config file, writing defaults...")
		if err := s.WriteConfigAs(".common.yaml"); err != nil {
			return nil, errors.WithStack(err)
		} else {
			return nil, errors.New("CONFIG WRITTEN TO CURRENT DIRECTORY- .common.yaml \nPLEASE MOVE CONFIG TO HOME DIRECTORY")
		}

	} else {
		log.Println("Using config file-->", zap.String("config_file", s.ConfigFileUsed()))
		if err := s.WriteConfig(); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return s, nil
}

func (s *Config) Annotate() map[string]string {
	settings := s.AllSettings()
	an := make(map[string]string)
	for k, v := range settings {
		if t, ok := v.(string); ok == true {
			an[k] = t
		}
	}
	return an
}
