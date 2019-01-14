//+build wireinject

package common

import (
	"context"
	"github.com/google/wire"
	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
	"gocloud.dev/gcp/gcpcloud"
	"gocloud.dev/mysql/cloudmysql"
)

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

func InjectApp(ctx context.Context) (*Application, func(), error) {
	wire.Build(
		ApplicationSet,
		NewConfig,
		gcpcloud.GCP,
		cloudmysql.Open,
		gcpSQLParams,
		gcpBucket,
	)
	return nil, nil, nil
}
