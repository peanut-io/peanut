package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/peanut-io/peanut/logger"
	"github.com/peanut-io/peanut/storage/core"
)

type Minio struct {
	client     *minio.Client
	bucketName string
}

func init() {
	core.RegisterStorage(core.Minio, NewMinio)
}

func NewMinio(cfg *core.Config) core.Storage {
	m := &Minio{}
	ctx := context.Background()
	if client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	}); err != nil {
		panic(fmt.Sprintf("failed to minio connect, endpoint: %s, error: %s", cfg.Endpoint, err.Error()))
	} else {
		m.client = client
	}
	if err := m.SetBucket(ctx, cfg.BucketName); err != nil {
		logger.Errorw(err.Error())
		return nil
	}
	return m
}

func (m *Minio) BucketName() string {
	return m.bucketName
}

func (m *Minio) SetBucket(ctx context.Context, bucketName string) error {
	if existed, err := m.client.BucketExists(ctx, bucketName); err != nil {
		logger.Infow("bucketName query error", "error", err)
		return err
	} else if !existed {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region:        "cn-north-1",
			ObjectLocking: true,
		})
		if err != nil {
			logger.Infow("minio failed to make buckName", "buckName", bucketName, "error", err)
			return err
		}
	}
	m.bucketName = bucketName
	return nil
}

func (m *Minio) Upload(ctx context.Context, src string, dst string) error {
	if _, err := m.client.FPutObject(ctx, m.bucketName, dst, src, minio.PutObjectOptions{}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("minio failed to upload %s to %s", src, dst))
	}
	return nil
}

func (m *Minio) Download(ctx context.Context, src string, dst string) error {
	if err := m.client.FGetObject(ctx, m.bucketName, src, dst, minio.GetObjectOptions{}); err != nil {
		errResponse := &minio.ErrorResponse{}
		if errors.As(err, errResponse) && errResponse.Code == "NoSuchKey" && strings.HasPrefix(src, "/") {
			src = src[1:]
			bucketName := m.bucketName
			if pos := strings.Index(src, "/"); pos > 0 {
				bucketName = src[0:pos]
				src = src[pos:]
			}
			if err = m.client.FGetObject(ctx, bucketName, src, dst, minio.GetObjectOptions{}); err != nil {
				logger.Errorw("minio failed to download object", "file", src, "error", err)
				return err
			}
		}
	}
	return nil
}

func (m *Minio) Read(ctx context.Context, src string) ([]byte, error) {
	obj, err := m.client.GetObject(ctx, m.bucketName, src, minio.GetObjectOptions{})
	defer func() {
		_ = obj.Close()
	}()
	if err != nil {
		errResponse := &minio.ErrorResponse{}
		if errors.As(err, errResponse) && errResponse.Code == "NoSuchKey" && strings.HasPrefix(src, "/") {
			src = src[1:]
			bucketName := m.bucketName
			if pos := strings.Index(src, "/"); pos > 0 {
				bucketName = src[0:pos]
				src = src[pos:]
			}
			if obj, err = m.client.GetObject(ctx, bucketName, src, minio.GetObjectOptions{}); err != nil {
				logger.Errorw("minio failed to get object", "file", src, "error", err)
				return nil, err
			}
		}
	}
	info, err := obj.Stat()
	if err != nil {
		errResponse := &minio.ErrorResponse{}
		if errors.As(err, errResponse) && errResponse.Code == "NoSuchKey" {
			return nil, errors.Errorf("minio failed to found the file %s", src)
		}
		return nil, err
	}
	buf := make([]byte, info.Size)
	if n, err := obj.Read(buf); (err != nil && err != io.EOF) || n != int(info.Size) {
		logger.Errorw("minio failed to read the content", "file", src, "error", err)
		return nil, errors.Errorf("minio failed to read the content, error %s", err.Error())
	}
	return buf, nil
}

func (m *Minio) Write(ctx context.Context, data []byte, dst string) error {
	_, err := m.client.PutObject(ctx, m.bucketName, dst, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	return err
}

func (m *Minio) WriteString(ctx context.Context, data string, dst string) error {
	_, err := m.client.PutObject(ctx, m.bucketName, dst, strings.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	return err
}
