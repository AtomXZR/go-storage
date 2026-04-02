package s3

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Options struct {
	UseSSL bool
	Region string
}

type Config struct {
	Bucket    string
	Endpoint  string
	AccessKey string
	SecretKey string

	UseSSL bool
	Region string
}

func New(conf Config) (*S3Storage, error) {
	client, err := minio.New(conf.Endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(conf.AccessKey, conf.SecretKey, ""),
		Secure:       conf.UseSSL,
		BucketLookup: minio.BucketLookupAuto,
		Region:       conf.Region,
	})

	if err != nil {
		return nil, s3ErrorToFSError(err)
	}

	return &S3Storage{
		client: client,
		bucket: conf.Bucket,
	}, nil
}
