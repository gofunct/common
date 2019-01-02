package gencmd

import (
	"github.com/gofunct/gogen/pkg/gogencmd"
)

// Option configures a command context.
type Option func(*Ctx)

func WithGogenCtx(gctx *gogencmd.Ctx) Option {
	return func(ctx *Ctx) {
		ctx.Ctx = gctx
	}
}

// WithCreateAppFunc specifies a dependencies initializer.
func WithCreateAppFunc(f CreateAppFunc) Option {
	return func(ctx *Ctx) {
		ctx.CreateAppFunc = f
	}
}
