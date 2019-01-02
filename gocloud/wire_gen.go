// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package gocloud

import (
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver/monitoredresource"
	"database/sql"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-sql-driver/mysql"
	"github.com/gofunct/common/gocloud/aws"
	"github.com/gofunct/common/gocloud/google"
	"go.opencensus.io/trace"
	"gocloud.dev/aws/rds"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/blob/s3blob"
	"gocloud.dev/gcp"
	"gocloud.dev/gcp/cloudsql"
	"gocloud.dev/mysql/cloudmysql"
	"gocloud.dev/mysql/rdsmysql"
	"gocloud.dev/requestlog"
	"gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/filevar"
	"gocloud.dev/runtimevar/runtimeconfigurator"
	"gocloud.dev/server"
	"gocloud.dev/server/sdserver"
	"gocloud.dev/server/xrayserver"
	"google.golang.org/genproto/googleapis/cloud/runtimeconfig/v1beta1"
	"net/http"
)

// Injectors from inject_aws.go:

func SetupAWS(ctx context.Context, flags *cliFlags) (*Application, func(), error) {
	ncsaLogger := xrayserver.NewRequestLogger()
	client := _wireClientValue
	certFetcher := &rds.CertFetcher{
		Client: client,
	}
	params := AwsSQLParams(flags)
	db, cleanup, err := rdsmysql.Open(ctx, certFetcher, params)
	if err != nil {
		return nil, nil, err
	}
	v, cleanup2 := AppHealthChecks(db)
	options := _wireOptionsValue
	sessionSession, err := session.NewSessionWithOptions(options)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	xRay := xrayserver.NewXRayClient(sessionSession)
	exporter, cleanup3, err := xrayserver.NewExporter(xRay)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	sampler := trace.AlwaysSample()
	defaultDriver := _wireDefaultDriverValue
	serverOptions := &server.Options{
		RequestLogger:         ncsaLogger,
		HealthChecks:          v,
		TraceExporter:         exporter,
		DefaultSamplingPolicy: sampler,
		Driver:                defaultDriver,
	}
	serverServer := server.New(serverOptions)
	bucket, err := AwsBucket(ctx, sessionSession, flags)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	variable, err := AwsRuntimeConfig(ctx, sessionSession, flags)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	application := NewApplication(serverServer, db, bucket, variable)
	return application, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

var (
	_wireClientValue        = http.DefaultClient
	_wireOptionsValue       = session.Options{}
	_wireDefaultDriverValue = &server.DefaultDriver{}
)

// Injectors from inject_gcp.go:

func SetupGCP(ctx context.Context, flags *cliFlags) (*Application, func(), error) {
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
	params := GcpSQLParams(projectID, flags)
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
	bucket, err := GcpBucket(ctx, flags, httpClient)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	runtimeConfigManagerClient, cleanup3, err := runtimeconfigurator.Dial(ctx, tokenSource)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	variable, cleanup4, err := GcpRuntimeConfig(ctx, runtimeConfigManagerClient, projectID, flags)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	application := NewApplication(serverServer, db, bucket, variable)
	return application, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// Injectors from inject_local.go:

func SetupLocal(ctx context.Context, flags *cliFlags) (*Application, func(), error) {
	logger := _wireLoggerValue
	db, err := DialLocalSQL(flags)
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
	bucket, err := LocalBucket(flags)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	variable, cleanup2, err := LocalRuntimeConfig(flags)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	application := NewApplication(serverServer, db, bucket, variable)
	return application, func() {
		cleanup2()
		cleanup()
	}, nil
}

var (
	_wireLoggerValue   = requestlog.Logger(nil)
	_wireExporterValue = trace.Exporter(nil)
)

// inject_aws.go:

// awsBucket is a Wire provider function that returns the S3 bucket based on the
// command-line flags.
func AwsBucket(ctx context.Context, cp client.ConfigProvider, flags *cliFlags) (*blob.Bucket, error) {
	return s3blob.OpenBucket(ctx, flags.bucket, cp, nil)
}

// awsSQLParams is a Wire provider function that returns the RDS SQL connection
// parameters based on the command-line flags. Other providers inside
// awscloud.AWS use the parameters to construct a *sql.DB.
func AwsSQLParams(flags *cliFlags) *rdsmysql.Params {
	return &rdsmysql.Params{
		Endpoint: flags.dbHost,
		Database: flags.dbName,
		User:     flags.dbUser,
		Password: flags.dbPassword,
	}
}

// awsMOTDVar is a Wire provider function that returns the Message of the Day
// variable from SSM Parameter Store.
func AwsRuntimeConfig(ctx context.Context, sess client.ConfigProvider, flags *cliFlags) (*runtimevar.Variable, error) {
	return aws.NewVariable(sess, flags.motdVar, runtimevar.StringDecoder, &aws.Options{
		WaitDuration: flags.motdVarWaitTime,
	})
}

// inject_gcp.go:

// gcpBucket is a Wire provider function that returns the GCS bucket based on
// the command-line flags.
func GcpBucket(ctx context.Context, flags *cliFlags, client2 *gcp.HTTPClient) (*blob.Bucket, error) {
	return gcsblob.OpenBucket(ctx, flags.bucket, client2, nil)
}

// gcpSQLParams is a Wire provider function that returns the Cloud SQL
// connection parameters based on the command-line flags. Other providers inside
// gcpcloud.GCP use the parameters to construct a *sql.DB.
func GcpSQLParams(id gcp.ProjectID, flags *cliFlags) *cloudmysql.Params {
	return &cloudmysql.Params{
		ProjectID: string(id),
		Region:    flags.cloudSQLRegion,
		Instance:  flags.dbHost,
		Database:  flags.dbName,
		User:      flags.dbUser,
		Password:  flags.dbPassword,
	}
}

// gcpMOTDVar is a Wire provider function that returns the Message of the Day
// variable from Runtime Configurator.
func GcpRuntimeConfig(ctx context.Context, client2 runtimeconfig.RuntimeConfigManagerClient, project gcp.ProjectID, flags *cliFlags) (*runtimevar.Variable, func(), error) {
	name := google.ResourceName{
		ProjectID: string(project),
		Config:    flags.runtimeConfigName,
		Variable:  flags.motdVar,
	}
	v, err := google.NewVariable(client2, name, runtimevar.StringDecoder, &google.Options{
		WaitDuration: flags.motdVarWaitTime,
	})
	if err != nil {
		return nil, nil, err
	}
	return v, func() { v.Close() }, nil
}

// inject_local.go:

// LocalBucket is a Wire provider function that returns a directory-based bucket
// based on the command-line flags.
func LocalBucket(flags *cliFlags) (*blob.Bucket, error) {
	return fileblob.OpenBucket(flags.bucket, nil)
}

// DialLocalSQL is a Wire provider function that connects to a MySQL database
// (usually on localhost).
func DialLocalSQL(flags *cliFlags) (*sql.DB, error) {
	cfg := &mysql.Config{
		Net:                  "tcp",
		Addr:                 flags.dbHost,
		DBName:               flags.dbName,
		User:                 flags.dbUser,
		Passwd:               flags.dbPassword,
		AllowNativePasswords: true,
	}
	return sql.Open("mysql", cfg.FormatDSN())
}

// LocalRuntimeVar is a Wire provider function that returns the Message of the
// Day variable based on a local file.
func LocalRuntimeConfig(flags *cliFlags) (*runtimevar.Variable, func(), error) {
	v, err := filevar.New(flags.motdVar, runtimevar.StringDecoder, &filevar.Options{
		WaitDuration: flags.motdVarWaitTime,
	})
	if err != nil {
		return nil, nil, err
	}
	return v, func() { v.Close() }, nil
}
