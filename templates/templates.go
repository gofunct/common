package templates

import (
	"github.com/iancoleman/strcase"
	"text/template"
)

var (
	FuncMap      = template.FuncMap{"ToCamel": strcase.ToCamel}
	TemplateMain = MustCreateTemplate("main", `package main

import (
	"fmt"
	"os"

	"github.com/izumin5210/clig/pkg/clib"

	"{{.Package}}/pkg/{{.Name}}/cmd"
)

const (
	appName = "{{.Name}}"
	version = "v0.0.1"
)

var (
	revision, buildDate string
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := cmd.NewDefault{{ToCamel .Name}}Command(clib.Path(wd), clib.Build{
		AppName:   appName,
		Version:   version,
		Revision:  revision,
		BuildDate: buildDate,
	})

	return cmd.Execute()
}
`)
)

func MustCreateTemplate(name, tmpl string) *template.Template {
	return template.Must(template.New(name).Funcs(FuncMap).Parse(tmpl))
}