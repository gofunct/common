package templates

import "text/template"

func ContextTemplate() *template.Template {
	return MustCreateTemplate("ctx", `package {{.Name}}

import (
	"github.com/izumin5210/clig/pkg/clib"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	{{- if .ViperEnabled}}
	"github.com/spf13/viper"
	{{- end}}
	"go.uber.org/zap"
	"k8s.io/utils/exec"
)

type Ctx struct {
	WorkingDir clib.Path
	IO         clib.IO
	FS         afero.Fs
	{{- if .ViperEnabled}}
	Viper      *viper.Viper
	{{- end}}
	Exec       exec.Interface

	Build  clib.Build
	Config *Config
}

func (c *Ctx) Init() error {
	{{- if .ViperEnabled}}
	c.Viper.SetFs(c.FS)

	var err error

	err = c.loadConfig()
	if err != nil {
		return errors.WithStack(err)
	}

	{{- end}}

	return nil
}
{{- if .ViperEnabled}}

func (c *Ctx) loadConfig() error {
	c.Viper.SetConfigName(c.Build.AppName)

	err := c.Viper.ReadInConfig()
	if err != nil {
		zap.L().Info("failed to find a config file", zap.Error(err))
		return nil
	}

	err = c.Viper.Unmarshal(c.Config)
	if err != nil {
		zap.L().Warn("failed to parse the config file", zap.Error(err))
		return errors.WithStack(err)
	}

	return nil
}
{{- end}}
`)
}
