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

func SetupGCP(ctx context.Context, cfg *Config) (*Application, func(), error) {
	wire.Build(
		gcpcloud.GCP,
		cloudmysql.Open,
		ApplicationSet,
		gcpBucket,
		gcpSQLParams,
	)
	return nil, nil, nil
}

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
