//+build wireinject

package common

import (
	"database/sql"
	"github.com/gofunct/common/ask"
	"github.com/gofunct/common/fs"
	"github.com/gofunct/common/log"
	"github.com/gofunct/common/render"
	"github.com/gofunct/common/router"
	"github.com/gofunct/iio"
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
	"go.opencensus.io/trace"
	"gocloud.dev/blob"
	"gocloud.dev/health"
	"gocloud.dev/health/sqlhealth"
	"gocloud.dev/server"
)

// newApplication creates a new Application struct based on the backends and the message
// of the day variable.
func NewApplication(srv *server.Server, db *sql.DB, bucket *blob.Bucket, config *Config, fs *fs.Service, q *ask.Service, r *render.Service, l *log.Service, i *iio.Service, rout *mux.Router) *Application {
	app := &Application{
		srv:      srv,
		db:       db,
		bucket:   bucket,
		Config:   config,
		FS:       fs,
		Q:        q,
		Renderer: r,
		L:        l,
		IO:       i,
		Router:   rout,
	}
	return app
}

var ApplicationSet = wire.NewSet(
	NewApplication,
	AppHealthChecks,
	trace.AlwaysSample,
	CommonSet,
)

var CommonSet = wire.NewSet(
	ask.Inject,
	fs.Inject,
	iio.Inject,
	log.InjectVerbose,
	router.Inject,
	render.Inject,
)

func NewConfig(set *pflag.FlagSet) (*Config, error) {
	c := &Config{}
	if err := c.Init(); err != nil {
		return nil, err
	}
	c.Bind(set)
	if err := c.BindPFlags(c.FlagSet); err != nil {
		return nil, err
	}

	return c, nil
}

func AppHealthChecks(db *sql.DB) ([]health.Checker, func()) {
	dbCheck := sqlhealth.New(db)
	list := []health.Checker{dbCheck}
	return list, func() {
		dbCheck.Stop()
	}
}
