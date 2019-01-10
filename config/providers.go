package config

import (
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

func New() (*Service, error) {
	s := &Service{}

	if err := rootInit(s); err != nil {
		return nil, err
	}

	if err := serviceInit(s); err != nil {
		return nil, err
	}

	return s, nil
}

var DefaultSet = wire.NewSet(
	New,
)

var serviceInit = func(s *Service) error {
	{
		s.Service = s.Sub("service")
		s.Service.SetConfigName("service")
		s.Service.AddConfigPath(os.Getenv("PWD"))
		s.Service.SetDefault("log-level", "verbose")
	}

	return nil
}

var rootInit = func(s *Service) error {
	s.Viper = viper.GetViper()
	if s.Viper == nil {
		s.Viper = viper.New()
	}

	s.AutomaticEnv()

	ex, err := os.Executable()
	if err != nil {
		return errors.WithStack(err)
	}
	s.SetDefault("executable", ex)
	gr, err := os.Getgroups()

	if err != nil {
		return errors.WithStack(err)
	}
	s.SetDefault("groups", gr)
	host, err := os.Hostname()
	if err != nil {
		return errors.WithStack(err)
	}

	s.SetDefault("env", os.Environ())
	s.SetDefault("uid", os.Getuid())
	s.SetDefault("args", os.Args)
	s.SetDefault("hostname", host)
	s.SetDefault("pid", os.Getpid())
	s.SetDefault("goarch", runtime.GOARCH)
	s.SetDefault("compiler", runtime.Compiler)
	s.SetDefault("version", runtime.Version())
	s.SetDefault("goos", runtime.GOOS)
	s.SetDefault("goroot", runtime.GOROOT())

	return nil
}
