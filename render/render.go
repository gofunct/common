package render

import (
	"bytes"
	"github.com/shiyanhui/hero"
	"os"
)

type Config struct {
	Source  string
	Dest    string
	PkgName string
}

func NewConfig() *Config {
	c := new(Config)
	c.Source = os.Getenv("RENDER_SOURCE")
	c.Dest = os.Getenv("RENDER_DEST")
	c.PkgName = os.Getenv("RENDER_PKG")
	return c
}

type Service struct {
	*Config
	Buffer *bytes.Buffer
}

func NewRenderer(config *Config) *Service {
	return &Service{
		Config: config,
		Buffer: hero.GetBuffer(),
	}
}

func (r *Service) Generate() {
	hero.Generate(r.Source, r.Dest, r.PkgName)
}
