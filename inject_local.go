//+build wireinject

package common

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"go.opencensus.io/trace"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/requestlog"
	"gocloud.dev/server"
)

func SetupLocal(ctx context.Context, cfg *Config) (*application, func(), error) {
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
