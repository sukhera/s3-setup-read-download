package main

import (
	"context"
	"fmt"
	"github.com/ahmed.sukhera/dls3obj/src/files"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	EnvVarBucketName = "BUCKET_NAME"
	EnvVarBucketPath = "BUCKET_PATH"
)

func main() {
	ctx := context.Background()
	loadEnv()

	bucketName := os.Getenv(EnvVarBucketName)
	bucketPath := os.Getenv(EnvVarBucketPath)

	awsConfig, err := config.LoadDefaultConfig(ctx)
	failOnErr(err)

	s3ApiClient := files.NewS3ApiClient(awsConfig)
	fileOpener := files.NewS3Store(s3ApiClient)

	data, err := files.DownloadFile(ctx, fileOpener, bucketName, bucketPath)
	failOnErr(err)

	fmt.Print(string(data))
}

// downloadFile stored in a AWS S3 bucket to a local file
func downloadFile(bucketName string, objectKey string, filename string) (err error) {
	// Create a file to download to
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION environment variables are used by LoadDefaultConfig
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return
	}

	// Create a S3 client using our configuration
	s3Client := s3.NewFromConfig(cfg)

	// Download the S3 object using the S3 manager object downloader
	downloader := manager.NewDownloader(s3Client)
	_, err = downloader.Download(context.TODO(), file, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	return
}

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

//loadEnv reads the content of the .env file.
// Note: Using .env files should be discouraged as they can cause a few problems.
// 1. Potential security risk
// 2. Source control issues for different environments (CI, Local, Staging, QA, Prod etc)
// Services like Vault, AWS Secrets Manager or GCP Secrets Manager should be preferred.
func loadEnv() {
	failOnErr(godotenv.Load(envPath))
}
