package s3

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage struct {
	client *s3.Client
	bucket string
}

type S3Config struct {
	URL    string
	Key    string
	Secret string
}

var _ FileStorage = (*Storage)(nil)

func NewStorage(bucket string, conf S3Config) *Storage {
	cfg := aws.Config{
		Region:      "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider(conf.Key, conf.Secret, ""),
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(conf.URL)
		o.UsePathStyle = true
	})
	ctx := context.Background()
	if err := ensureBucketExists(ctx, client, bucket); err != nil {
		log.Fatalf("Error creating s3 bucket: %v\n", err)
	}
	return &Storage{
		client: client,
		bucket: bucket,
	}
}
