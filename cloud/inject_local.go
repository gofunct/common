//+build wireinject

package cloud

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"go.opencensus.io/trace"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/requestlog"
	"gocloud.dev/runtimevar"
	"gocloud.dev/runtimevar/filevar"
	"gocloud.dev/server"
)

// This file wires the generic interfaces up to local implementations. It won't
// be directly included in the final binary, since it includes a Wire injector
// template function (setupLocal), but the declarations will be copied into
// wire_gen.go when Wire is run.

// setupLocal is a Wire injector function that sets up the application using
// local implementations.
func SetupLocal(ctx context.Context, flags *cliFlags) (*application, func(), error) {
	// This will be filled in by Wire with providers from the provider sets in
	// wire.Build.
	wire.Build(
		wire.InterfaceValue(new(requestlog.Logger), requestlog.Logger(nil)),
		wire.InterfaceValue(new(trace.Exporter), trace.Exporter(nil)),
		server.Set,
		ApplicationSet,
		dialLocalSQL,
		localBucket,
	)
	return nil, nil, nil
}

// localBucket is a Wire provider function that returns a directory-based bucket
// based on the command-line flags.
func localBucket(flags *cliFlags) (*blob.Bucket, error) {
	return fileblob.OpenBucket(flags.bucket, nil)
}

// dialLocalSQL is a Wire provider function that connects to a MySQL database
// (usually on localhost).
func dialLocalSQL(flags *cliFlags) (*sql.DB, error) {
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

// localRuntimeVar is a Wire provider function that returns the Message of the
// Day variable based on a local file.
func localRuntimeVar(flags *cliFlags) (*runtimevar.Variable, func(), error) {
	v, err := filevar.New(flags.motdVar, runtimevar.StringDecoder, &filevar.Options{
		WaitDuration: flags.motdVarWaitTime,
	})
	if err != nil {
		return nil, nil, err
	}
	return v, func() { v.Close() }, nil
}
