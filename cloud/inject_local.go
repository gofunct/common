//+build wireinject

package cloud

import (
	"context"
	"database/sql"
	"github.com/gofunct/common/config"

	"github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"go.opencensus.io/trace"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/requestlog"
	"gocloud.dev/server"
)

func SetupLocal(ctx context.Context, cfg *config.Service) (*application, func(), error) {
	wire.Build(
		wire.InterfaceValue(new(requestlog.Logger), requestlog.Logger(nil)),
		wire.InterfaceValue(new(trace.Exporter), trace.Exporter(nil)),
		server.Set,
		applicationSet,
		dialLocalSQL,
		localBucket,
	)
	return nil, nil, nil
}

func localBucket(cfg *config.Service) (*blob.Bucket, error) {
	return fileblob.OpenBucket(cfg.Bucket, nil)
}

func dialLocalSQL(c *config.Service) (*sql.DB, error) {
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
