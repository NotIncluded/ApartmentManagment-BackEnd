package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client     *minio.Client
	BucketName string
}

func NewMinioClient(endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	// Ensure bucket exists
	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}
	return &MinioClient{Client: client, BucketName: bucketName}, nil
}

func (m *MinioClient) UploadFile(ctx context.Context, objectName, filePath, contentType string) (string, error) {
	info, err := m.Client.FPutObject(ctx, m.BucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", m.Client.EndpointURL(), m.BucketName, info.Key), nil
}
