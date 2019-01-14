package common

import (
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver/monitoredresource"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/gofunct/common/ask"
	"github.com/gofunct/common/fs"
	"github.com/gofunct/common/log"
	"github.com/gofunct/common/render"
	"github.com/gofunct/common/router"
	"github.com/gofunct/iio"
	"go.opencensus.io/trace"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
	"gocloud.dev/gcp/cloudsql"
	"gocloud.dev/mysql/cloudmysql"
	"gocloud.dev/requestlog"
	"gocloud.dev/server"
	"gocloud.dev/server/sdserver"
)

// Injectors from inject_gcp.go:

func SetupGCP(ctx context.Context, cfg *Config) (*Application, func(), error) {
	stackdriverLogger := sdserver.NewRequestLogger()
	roundTripper := gcp.DefaultTransport()
	credentials, err := gcp.DefaultCredentials(ctx)
	if err != nil {
		return nil, nil, err
	}
	tokenSource := gcp.CredentialsTokenSource(credentials)
	httpClient, err := gcp.NewHTTPClient(roundTripper, tokenSource)
	if err != nil {
		return nil, nil, err
	}
	remoteCertSource := cloudsql.NewCertSource(httpClient)
	projectID, err := gcp.DefaultProjectID(credentials)
	if err != nil {
		return nil, nil, err
	}
	params := gcpSQLParams(projectID, cfg)
	db, err := cloudmysql.Open(ctx, remoteCertSource, params)
	if err != nil {
		return nil, nil, err
	}
	v, cleanup := AppHealthChecks(db)
	monitoredresourceInterface := monitoredresource.Autodetect()
	exporter, cleanup2, err := sdserver.NewExporter(projectID, tokenSource, monitoredresourceInterface)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	sampler := trace.AlwaysSample()
	defaultDriver := _wireDefaultDriverValue
	options := &server.Options{
		RequestLogger:         stackdriverLogger,
		HealthChecks:          v,
		TraceExporter:         exporter,
		DefaultSamplingPolicy: sampler,
		Driver:                defaultDriver,
	}
	serverServer := server.New(options)
	bucket, err := gcpBucket(ctx, cfg, httpClient)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	service, err := fs.Inject()
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	askService, err := ask.Inject()
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	renderService := render.Inject()
	logService, err := log.InjectVerbose()
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	iioService, err := iio.Inject()
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	muxRouter := router.Inject()
	application := NewApplication(serverServer, db, bucket, cfg, service, askService, renderService, logService, iioService, muxRouter)
	return application, func() {
		cleanup2()
		cleanup()
	}, nil
}

var (
	_wireDefaultDriverValue = &server.DefaultDriver{}
)

// Injectors from inject_local.go:

func SetupLocal(ctx context.Context, cfg *Config) (*Application, func(), error) {
	logger := _wireLoggerValue
	db, err := dialLocalSQL(cfg)
	if err != nil {
		return nil, nil, err
	}
	v, cleanup := AppHealthChecks(db)
	exporter := _wireExporterValue
	sampler := trace.AlwaysSample()
	defaultDriver := _wireDefaultDriverValue
	options := &server.Options{
		RequestLogger:         logger,
		HealthChecks:          v,
		TraceExporter:         exporter,
		DefaultSamplingPolicy: sampler,
		Driver:                defaultDriver,
	}
	serverServer := server.New(options)
	bucket, err := localBucket(cfg)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	service, err := fs.Inject()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	askService, err := ask.Inject()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	renderService := render.Inject()
	logService, err := log.InjectVerbose()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	iioService, err := iio.Inject()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	muxRouter := router.Inject()
	application := NewApplication(serverServer, db, bucket, cfg, service, askService, renderService, logService, iioService, muxRouter)
	return application, func() {
		cleanup()
	}, nil
}

var (
	_wireLoggerValue   = requestlog.Logger(nil)
	_wireExporterValue = trace.Exporter(nil)
)

// inject_gcp.go:

func gcpBucket(ctx context.Context, cfg *Config, client *gcp.HTTPClient) (*blob.Bucket, error) {
	return gcsblob.OpenBucket(ctx, cfg.Bucket, client, nil)
}

func gcpSQLParams(id gcp.ProjectID, cfg *Config) *cloudmysql.Params {
	return &cloudmysql.Params{
		ProjectID: string(id),
		Region:    cfg.CloudSqlRegion,
		Instance:  cfg.DbHost,
		Database:  cfg.DbName,
		User:      cfg.DbUser,
		Password:  cfg.DbPassword,
	}
}

// inject_local.go:

func localBucket(cfg *Config) (*blob.Bucket, error) {
	return fileblob.OpenBucket(cfg.Bucket, nil)
}

func dialLocalSQL(c *Config) (*sql.DB, error) {
	cfg := &mysql.Config{
		Net:                  "tcp",
		Addr:                 c.DbHost,
		DBName:               c.DbName,
		User:                 c.DbUser,
		Passwd:               c.DbPassword,
		AllowNativePasswords: true,
	}
	return sql.Open("mysql", cfg.FormatDSN())
}
