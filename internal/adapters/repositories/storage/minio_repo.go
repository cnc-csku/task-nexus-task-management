package storage

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/minio/minio-go/v7"
)

type MinioRepositoryImpl struct {
	client *minio.Client
	config *config.Config
}

func NewMinioRepository(client *minio.Client, config *config.Config) repositories.MinioRepository {
	return &MinioRepositoryImpl{
		client: client,
		config: config,
	}
}

func (m *MinioRepositoryImpl) Upload(ctx context.Context, key string, object *os.File, contentType string) error {
	objectStat, err := object.Stat()
	if err != nil {
		return err
	}

	size := objectStat.Size()
	if size > m.config.MinioClient.FileUploadSizeLimitMB<<20 {
		return errors.New("file size limit exceeded")
	}

	_, err = m.client.PutObject(
		ctx,
		m.config.MinioClient.BucketName,
		key,
		object,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *MinioRepositoryImpl) GeneratePutPresignedURL(ctx context.Context, key string) (string, error) {
	presignedURL, err := m.client.PresignedPutObject(
		ctx,
		m.config.MinioClient.BucketName,
		key,
		time.Duration(m.config.MinioClient.PresignedURLExpirySec)*time.Second,
	)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (m *MinioRepositoryImpl) GetFullURL(key string) string {
	return m.config.MinioClient.Endpoint + "/" + m.config.MinioClient.BucketName + "/" + key
}
