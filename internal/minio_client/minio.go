package minio_client

import (
	"robot_scheduler/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var client *minio.Client

// Init 初始化 MinIO 客户端
func Init(cfg *config.MinioConfig) error {
	if cfg == nil || !cfg.Enabled {
		return nil
	}

	c, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return err
	}

	client = c
	return nil
}

// Client 获取全局 MinIO 客户端
func Client() *minio.Client {
	return client
}
