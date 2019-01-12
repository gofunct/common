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

type Renderer struct {
	*Config
	Buffer *bytes.Buffer
}

func NewRenderer(config *Config) *Renderer {
	return &Renderer{
		Config: config,
		Buffer: hero.GetBuffer(),
	}
}

func (r *Renderer) Generate() {
	hero.Generate(r.Source, r.Dest, r.PkgName)
}
