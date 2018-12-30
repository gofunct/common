package cobra

import (
	"github.com/gofunct/common/build"
	"github.com/gofunct/common/files"
	"github.com/gofunct/common/io"
	"github.com/gofunct/common/proto/protoc"
	"github.com/google/wire"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"k8s.io/utils/exec"
	"os"
	"path/filepath"
	"github.com/spf13/viper"
	"github.com/gofunct/common/errors"
)

func defaultCtx() *Context {
	return &Context{
		Ctx: &Ctx{},
	}
}

// Ctx defines a context of a generator.
type Context struct {
	*Ctx

	CreateAppFunc CreateAppFunc
}

func (c *Ctx) apply(opts []Option) {
	for _, f := range opts {
		f(c)
	}
}

// CreateApp initializes dependencies.
func (c *Ctx) CreateApp(cmd *Command) (*App, error) {
	f := c.CreateAppFunc
	if c.CreateAppFunc == nil {
		f = newApp
	}
	app, err := f(cmd)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return app, nil
}


// Ctx contains the runtime context of grpai.
type Ctx struct {
	FS     afero.Fs
	Viper  *viper.Viper
	Execer exec.Interface
	IO     io.IO

	RootDir   files.RootDir
	insideApp bool

	Config       Config
	Build        build.Build
	ProtocConfig protoc.Config
}

// Config stores general setting params and provides accessors for them.
type Config struct {
	Package string
	Grapi   struct {
		ServerDir string
	}
}

// Init initializes the runtime context.
func (c *Ctx) Init() error {
	if c.RootDir.String() == "" {
		dir, _ := os.Getwd()
		c.RootDir = files.RootDir{files.Path(dir)}
	}

	if c.IO == nil {
		c.IO = io.Stdio()
	}

	if c.FS == nil {
		c.FS = afero.NewOsFs()
	}

	if c.Viper == nil {
		c.Viper = viper.New()
	}

	c.Viper.SetFs(c.FS)

	if c.Execer == nil {
		c.Execer = exec.New()
	}

	if c.Build.AppName == "" {
		c.Build.AppName = "grapi"
	}

	return errors.WithStack(c.loadConfig())
}

func (c *Ctx) loadConfig() error {
	c.Viper.SetConfigName(c.Build.AppName)
	for dir := c.RootDir.String(); dir != "/"; dir = filepath.Dir(dir) {
		c.Viper.AddConfigPath(dir)
	}

	err := c.Viper.ReadInConfig()
	if err != nil {
		zap.L().Info("failed to find config file", zap.Error(err))
		return nil
	}

	c.insideApp = true
	c.RootDir = files.RootDir{files.Path(filepath.Dir(c.Viper.ConfigFileUsed()))}

	err = c.Viper.Unmarshal(&c.Config)
	if err != nil {
		zap.L().Warn("failed to parse config", zap.Error(err))
		return errors.WithStack(err)
	}

	err = c.Viper.UnmarshalKey("protoc", &c.ProtocConfig)
	if err != nil {
		zap.L().Warn("failed to parse protoc config", zap.Error(err))
		return errors.WithStack(err)
	}

	return nil
}

// IsInsideApp returns true if the current working directory is inside a grapi project.
func (c *Ctx) IsInsideApp() bool {
	return c.insideApp
}

// CtxSet is a provider set that includes modules contained in Ctx.
var CtxSet = wire.NewSet(
	ProvideFS,
	ProvideViper,
	ProvideExecer,
	ProvideIO,
	ProvideRootDir,
	ProvideConfig,
	ProvideBuildConfig,
	ProvideProtocConfig,
)

func ProvideFS(c *Ctx) afero.Fs                 { return c.FS }
func ProvideViper(c *Ctx) *viper.Viper          { return c.Viper }
func ProvideExecer(c *Ctx) exec.Interface       { return c.Execer }
func ProvideIO(c *Ctx) io.IO                  { return c.IO }
func ProvideRootDir(c *Ctx) files.RootDir         { return c.RootDir }
func ProvideConfig(c *Ctx) *Config              { return &c.Config }
func ProvideBuildConfig(c *Ctx) *build.Build     { return &c.Build }
func ProvideProtocConfig(c *Ctx) *protoc.Config { return &c.ProtocConfig }
