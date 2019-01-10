package config

import (
	"encoding/json"
	"github.com/spf13/viper"
)

type API struct {
	object interface{}
	V      *viper.Viper
	provider string
    endpoint string
    url      string
    secret   string
}

func (a API) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.object)
}

func (a API) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, a.object)
}
func (a API) SetObject(i interface{}) {
	a.object = i
}

func (a API) MergeConfigMap(cfg map[string]interface{}) error {
	return a.V.MergeConfigMap(cfg)
}

func (a API) AllSettings() map[string]interface{} {
	return a.V.AllSettings()
}

func (a API) Set(k string, V interface{}) {
	a.V.Set(k, V)
}

func (a API) GetString(k string) string {
	return a.V.GetString(k)
}

func (a API) GetStringSlice(k string) []string {
	return a.V.GetStringSlice(k)
}

func (a API) BindEnv(k string) {
	a.V.BindEnv(k)
}

func (a API) GetBool(k string) bool {
	return a.V.GetBool(k)
}

func (a API) Unmarshal() error {
	return a.V.Unmarshal(a.object)
}

func (a API) Debug() {
	a.V.Debug()
}
func (a API) ReadInConfig() error {
	return a.V.ReadInConfig()
}

func (a API) ReadRemote() error {
	return a.V.ReadRemoteConfig()
}

func (a API) GetObject() interface{} {
	return a.object
}
func (r *API) Provider() string {
	return r.provider
}

func (r *API) Endpoint() string {
	return r.endpoint
}

func (r *API) Path() string {
	return r.endpoint
}

func (r *API) SecretKeyring() string {
	return r.endpoint
}

func (r *API) SetProvider(s string) {
	r.provider = s
}

func (r *API) SetEndpoint(s string) {
	r.provider = s
}

func (r *API) SetPath(s string) {
	r.provider = s
}

func (r *API) SetSecretKeyring(s string) {
	r.provider = s
}
