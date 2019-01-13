//+build wireinject

package cloud

import (
	"context"
	"github.com/gofunct/common/config"

	"github.com/google/wire"
	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
	"gocloud.dev/gcp/gcpcloud"
	"gocloud.dev/mysql/cloudmysql"
)

func SetupGCP(ctx context.Context, cfg *config.Service) (*application, func(), error) {
	wire.Build(
		gcpcloud.GCP,
		cloudmysql.Open,
		applicationSet,
		gcpBucket,
		gcpSQLParams,
	)
	return nil, nil, nil
}

func gcpBucket(ctx context.Context, cfg *config.Service, client *gcp.HTTPClient) (*blob.Bucket, error) {
	return gcsblob.OpenBucket(ctx, cfg.Bucket, client, nil)
}

func gcpSQLParams(id gcp.ProjectID, cfg *config.Service) *cloudmysql.Params {
	return &cloudmysql.Params{
		ProjectID: string(id),
		Region:    cfg.CloudSqlRegion,
		Instance:  cfg.DbHost,
		Database:  cfg.DbName,
		User:      cfg.DbUser,
		Password:  cfg.DbPassword,
	}
}
