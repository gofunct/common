package common

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"text/template"
)

// Wrapper can execute protoc commands for current project's proto files.
type Service interface {
	Exec(context.Context) error
}

type Handler func(ctx context.Context, request interface{}) (response interface{}, err error)

type Middleware func(Handler) Handler

type Runner func(*cobra.Command, []string) error

func (r *Runtime) Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}
func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
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
