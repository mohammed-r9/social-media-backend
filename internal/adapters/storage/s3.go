package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	client *s3.Client
	bucket string
}

type S3Config struct {
	URL    string
	Key    string
	Secret string
}

var _ Storage = (*S3)(nil)

func NewS3(bucket string, conf S3Config) *S3 {
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
	return &S3{
		client: client,
		bucket: bucket,
	}
}

func (s *S3) Upload(
	ctx context.Context,
	key string,
	body io.Reader,
	contentType ContentType,
) error {
	var seekableBody io.ReadSeeker

	switch v := body.(type) {
	case io.ReadSeeker:
		seekableBody = v
	default:
		b, err := io.ReadAll(body)
		if err != nil {
			return err
		}
		seekableBody = bytes.NewReader(b)
	}
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        seekableBody,
		ContentType: aws.String(contentType.String()),
	})

	return err
}

func (s *S3) GetURL(ctx context.Context, key string) (string, error) {
	presignClient := s3.NewPresignClient(s.client)
	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (s *S3) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

func ensureBucketExists(ctx context.Context, client *s3.Client, bucket string) error {
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err == nil {
		return nil
	}

	_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	return nil
}
