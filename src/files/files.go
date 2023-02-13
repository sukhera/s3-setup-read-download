package files

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type (
	// Store is a file store implementation.
	// Ideally this is located ins a store.go file in the package
	Store interface {
		FileOpener
	}

	// FileOpener is an interface that enables opening files.
	// Note: Using bucket and key for brevity of implementation.
	// If we need to support opening local files the method signature would change.
	FileOpener interface {
		OpenFile(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	}
)

func DownloadFile(
	ctx context.Context,
	fileOpener FileOpener,
	bucketName, bucketPath string,
) ([]byte, error) {
	if fileOpener == nil {
		return nil, fmt.Errorf("fileOpener implentation is required")
	}

	if bucketName == "" {
		return nil, fmt.Errorf("bucketName is required")
	}

	if bucketPath == "" {
		return nil, fmt.Errorf("bucketPath is required")
	}

	rc, err := fileOpener.OpenFile(ctx, bucketName, bucketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: Error %w", err)
	}
	defer rc.Close()

	return ioutil.ReadAll(rc)
}

/*****************************************************************************************/

type (

	// S3Api is a local interface that's compatible with the AWS S3 SDK.
	// Its purpose is to ease testing and mocking calls to the AWS api.
	// Ideally this should be in a separate s3.go file in the files package. :)
	S3Api interface {
		GetObject(
			ctx context.Context,
			params *s3.GetObjectInput,
			optFns ...func(*s3.Options),
		) (*s3.GetObjectOutput, error)
	}

	S3Store struct {
		s3Client S3Api
	}
)

// Compile time check that S3Store implements the Store interface
var _ Store = (*S3Store)(nil)

func NewS3Store(s3ApiClient S3Api) *S3Store {
	return &S3Store{
		s3Client: s3ApiClient,
	}
}

// OpenFile gets a file from S3 and returns a readCloser or an error. flag and mode are not used.
func (s *S3Store) OpenFile(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	out, err := s.s3Client.GetObject(ctx, params)
	if err == nil {
		return out.Body, nil
	}
	return nil, err
}

/*****************************************************************************************/

// s3ApiClient is an adapter for the aws s3 client.
type s3ApiClient struct {
	client *s3.Client
}

// Compile time check that s3ApiClient implements S3Api
var _ S3Api = (*s3ApiClient)(nil)

// NewS3ApiClient creates a new S3 client.
func NewS3ApiClient(cfg aws.Config) *s3ApiClient {
	return &s3ApiClient{client: s3.NewFromConfig(cfg)}
}

func (s *s3ApiClient) GetObject(
	ctx context.Context,
	params *s3.GetObjectInput,
	optFns ...func(*s3.Options),
) (*s3.GetObjectOutput, error) {
	return s.client.GetObject(ctx, params, optFns...)
}
