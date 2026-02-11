package minio_client

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"robot_scheduler/internal/config"

	"github.com/minio/minio-go/v7"
)

func initConfig(t *testing.T) {
	t.Helper()

	cfgPath := filepath.Join("..", "..", "configs", "config.yaml")
	// 允许通过环境变量覆盖配置路径
	if env := os.Getenv("ROBOT_SCHEDULER_CONFIG"); env != "" {
		cfgPath = env
	}

	if err := config.Init(cfgPath); err != nil {
		t.Fatalf("初始化配置失败: %v", err)
	}
}

func TestMinioConnection(t *testing.T) {
	initConfig(t)

	cfg := config.Get()
	if cfg.Minio == nil {
		t.Fatal("配置中未找到 minio 配置")
	}
	if !cfg.Minio.Enabled {
		t.Skip("minio.enabled = false，跳过连接测试")
	}

	if err := Init(cfg.Minio); err != nil {
		t.Fatalf("初始化 MinIO 客户端失败: %v", err)
	}

	c := Client()
	if c == nil {
		t.Fatal("MinIO 客户端为空")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := c.BucketExists(ctx, cfg.Minio.BucketName)
	if err != nil {
		t.Fatalf("检查 Bucket 失败: %v", err)
	}

	if !exists {
		t.Logf("Bucket %q 不存在，尝试创建", cfg.Minio.BucketName)
		if err := c.MakeBucket(ctx, cfg.Minio.BucketName, minio.MakeBucketOptions{Region: cfg.Minio.Region}); err != nil {
			t.Fatalf("创建 Bucket 失败: %v", err)
		}
	}

	t.Logf("MinIO 连接正常，Bucket=%s", cfg.Minio.BucketName)
}

func TestMinioUploadFile(t *testing.T) {
	initConfig(t)

	cfg := config.Get()
	if cfg.Minio == nil {
		t.Fatal("配置中未找到 minio 配置")
	}
	if !cfg.Minio.Enabled {
		t.Skip("minio.enabled = false，跳过上传测试")
	}

	if err := Init(cfg.Minio); err != nil {
		t.Fatalf("初始化 MinIO 客户端失败: %v", err)
	}

	c := Client()
	if c == nil {
		t.Fatal("MinIO 客户端为空")
	}

	// 准备一个临时测试文件内容
	tmpFile, err := os.CreateTemp("", "minio-upload-test-*.txt")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	content := []byte("this is a minio upload go test\n")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("写入临时文件失败: %v", err)
	}

	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatalf("重置文件指针失败: %v", err)
	}

	stat, err := tmpFile.Stat()
	if err != nil {
		t.Fatalf("获取临时文件信息失败: %v", err)
	}

	objectName := "test/" + stat.Name()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	info, err := c.PutObject(ctx, cfg.Minio.BucketName, objectName, tmpFile, stat.Size(), minio.PutObjectOptions{
		ContentType: "text/plain",
	})
	if err != nil {
		t.Fatalf("上传对象失败: %v", err)
	}

	t.Logf("上传成功: bucket=%s, object=%s, size=%d, etag=%s",
		cfg.Minio.BucketName, objectName, info.Size, info.ETag)
}
