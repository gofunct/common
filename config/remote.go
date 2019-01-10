package config

import "github.com/spf13/viper"

type Provider interface {
	viper.RemoteProvider
	SetPath(string)
	SetProvider(string)
	SetEndpoint(string)
	SetSecretKeyring(string)
}

type RemoteConfig struct {
	provider string
	endpoint string
	url      string
	secret   string
}

func (r *RemoteConfig) Provider() string {
	return r.provider
}

func (r *RemoteConfig) Endpoint() string {
	return r.endpoint
}

func (r *RemoteConfig) Path() string {
	return r.endpoint
}

func (r *RemoteConfig) SecretKeyring() string {
	return r.endpoint
}

func (r *RemoteConfig) SetProvider(s string) {
	r.provider = s
}

func (r *RemoteConfig) SetEndpoint(s string) {
	r.provider = s
}

func (r *RemoteConfig) SetPath(s string) {
	r.provider = s
}

func (r *RemoteConfig) SetSecretKeyring(s string) {
	r.provider = s
}
