package cobra

// Option configures a command context.
type Option func(*Context)

// WithGrapiCtx specifies a grapi command context.
func WithCtx(gctx *Ctx) Option {
	return func(ctx *Context) {
		ctx = gctx
	}
}

// WithCreateAppFunc specifies a dependencies initializer.
func WithCreateAppFunc(f CreateAppFunc) Option {
	return func(ctx *Context) {
		CreateAppFunc = f
	}
}
