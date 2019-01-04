package files

import (
	"bytes"
	"github.com/gofunct/common/errors"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"github.com/serenize/snaker"
	"strings"
	"text/template"
)

var (
	FuncMap = template.FuncMap{"ToCamel": strcase.ToCamel}
)

type FuncMapper struct {
	Strings String
}

type String struct {
	Camel struct {
		Plural   string
		Singular string
	}
	CamelLower struct {
		Plural   string
		Singular string
	}
	Snake struct {
		Plural   string
		Singular string
	}
}

func Inflect(in string) (out String) {
	out.Camel.Singular = inflection.Singular(snaker.SnakeToCamel(in))
	out.Camel.Plural = inflection.Plural(out.Camel.Singular)
	out.CamelLower.Singular = strings.ToLower(string(out.Camel.Singular[0])) + out.Camel.Singular[1:]
	out.CamelLower.Plural = strings.ToLower(string(out.Camel.Plural[0])) + out.Camel.Plural[1:]
	out.Snake.Singular = snaker.CamelToSnake(out.Camel.Singular)
	out.Snake.Plural = snaker.CamelToSnake(out.Camel.Plural)
	return
}

func MustCreateTemplate(name, tmpl string) *template.Template {
	return template.Must(template.New(name).Funcs(FuncMap).Parse(tmpl))
}

// TemplateString is a compilable string with text/template package
type TemplateString string

// Compile generates textual output applied a parsed template to the specified values
func (s TemplateString) Compile(v interface{}) (string, error) {
	tmpl, err := template.New("").Parse(string(s))
	if err != nil {
		return string(s), errors.Wrapf(err, "failed to parse a template: %q", string(s))
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, v)
	if err != nil {
		return string(s), errors.Wrapf(err, "failed to execute a template: %q", string(s))
	}
	return string(buf.Bytes()), nil
}
