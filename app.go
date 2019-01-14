package common

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/mattn/go-colorable"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shiyanhui/hero"
	"github.com/spf13/afero"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"gocloud.dev/blob"
	"gocloud.dev/health"
	"gocloud.dev/health/sqlhealth"
	"gocloud.dev/server"
	"net/http"
	"net/http/pprof"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Application struct {
	Server  *server.Server
	Db      *sql.DB
	Bucket  *blob.Bucket
	Config  *Config
	Os      *afero.Afero
	Z       *zap.Logger
	IO      *IO
	Router  *mux.Router
	RunFunc func(context.Context, *Application) error
}

var ApplicationSet = wire.NewSet(
	NewApplication,
	AppHealthChecks,
	trace.AlwaysSample,
)

func NewApplication(srv *server.Server, db *sql.DB, bucket *blob.Bucket, cfg *Config) *Application {
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
	l.With(
		zap.String("user", os.Getenv("USER")),
		zap.Int("cpus", runtime.NumCPU()),
		zap.Int("routines", runtime.NumGoroutine()))

	a := &afero.Afero{
		Fs: afero.NewOsFs(),
	}
	router := mux.NewRouter()
	router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	router.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	router.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	app := &Application{
		Server: srv,
		Db:     db,
		Bucket: bucket,
		Config: cfg,
		Os:     a,
		Z:      l,
		IO: &IO{
			InR:  os.Stdin,
			OutW: colorable.NewColorableStdout(),
			ErrW: colorable.NewColorableStderr(),
		},
		Router:  router,
		RunFunc: nil,
	}
	return app
}

func (a *Application) Generate() {
	hero.Generate(a.Config.GenSource, a.Config.GenDest, a.Config.GenPkgName)
}

func (a *Application) Shell(args ...string) (stdout string, err error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	stdoutb, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("running %v: %v", cmd.Args, err)
	}
	return strings.TrimSpace(string(stdoutb)), nil
}

func (a *Application) Execute(ctx context.Context) error {
	return a.RunFunc(ctx, a)
}

func AppHealthChecks(db *sql.DB) ([]health.Checker, func()) {
	dbCheck := sqlhealth.New(db)
	list := []health.Checker{dbCheck}
	return list, func() {
		dbCheck.Stop()
	}
}
