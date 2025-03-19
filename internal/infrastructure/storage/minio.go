package storage

import (
	"context"
	"log"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOClient(ctx context.Context, cfg *config.Config) *minio.Client {
	client, err := minio.New(
		cfg.MinioClient.Endpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				cfg.MinioClient.AccessKeyID,
				cfg.MinioClient.SecretAccessKey,
				"",
			),
			Secure: cfg.MinioClient.UseSSL,
		})
	if err != nil {
		log.Fatalln("🚫 Cannot create MinIO client | ", err)
	}

	// test connection
	exists, err := client.BucketExists(ctx, cfg.MinioClient.BucketName)
	if err != nil {
		log.Fatalln("🚫 Cannot connect to MinIO | ", err)
	} else if !exists {
		log.Fatalf("🚫 Bucket %s does not exist", cfg.MinioClient.BucketName)
	} else {
		log.Println("✅ Connected to MinIO | Bucket:", cfg.MinioClient.BucketName)
	}

	return client
}
