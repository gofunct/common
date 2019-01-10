package config

import (
	"encoding/json"
	"github.com/spf13/viper"
)

type API struct {
	object interface{}
	*viper.Viper
	Flags []FlagValue
}

func (e API) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.object)
}

func (e API) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, e.object)
}

func (e API) AddFlag(f FlagValue) {
	e.Flags = append(e.Flags, f)
}

func (e API) SetObject(i interface{}) {
	e.object = i
}

func (a *API) VisitAll(fn func(FlagValue)) {
	for _, k := range a.Flags {
		fn(k)
	}
}
