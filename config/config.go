package config

import (
	"github.com/spf13/viper"
	"os/user"
)

type Service struct {
	Service *viper.Viper
	*viper.Viper
}

func (s *Service) Provider() string {
	return s.Service.GetString("remote.provider")
}

func (s *Service) Endpoint() string {
	return s.Service.GetString("remote.endpoint")
}

func (s *Service) Path() string {
	return s.Service.GetString("remote.path")
}

func (s *Service) SecretKeyring() string {
	return s.Service.GetString("remote.secret-key-ring")
}

func (s *Service) SetProvider(p string) {
	s.Service.SetDefault("remote.provider", p)
}

func (s *Service) SetEndpoint(p string) {
	s.Service.SetDefault("remote.endpoint", p)
}

func (s *Service) SetPath(p string) {
	s.Service.SetDefault("remote.path", p)
}

func (s *Service) SetSecretKeyring(p string) {
	s.Service.SetDefault("remote.secret", p)
}

func (s *Service) GetCurrentUser() string {
	usr, _ := user.Current()
	s.SetDefault("user", usr)
	return usr.Username
}
