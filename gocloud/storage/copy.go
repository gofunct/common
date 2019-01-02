package storage

import (
	"context"
	"gocloud.dev/blob"

	 "io"
	"os"
)

func CopyFileFromBucket(b *blob.Bucket, ctx context.Context, file string) error {
	r, err := b.NewReader(ctx, file, nil)
	if err != nil {
		return err
	}
	defer r.Close()

	if _, err := io.Copy(os.Stdout, r); err != nil {
		return err
	}
	return nil
}
