package oss

import (
	"context"
	"path"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	bucketName = "common"
)

type minioImpl struct {
	cli  *minio.Client
	once sync.Once
}

func NewMinio() Storage {
	if Impl != nil {
		return Impl
	}
	Impl = &minioImpl{once: sync.Once{}}
	return Impl
}

func (m *minioImpl) init() {
	m.once.Do(func() {
		var err error
		m.cli, err = minio.New("10.0.0.4:9000", &minio.Options{
			Creds:  credentials.NewStaticV4("TsP6KT9JeclMRzi6eJfm", "uGAsMuOB0Njwn1mCJRc0HrrB6Vfk5KfUWBKGAeKO", ""),
			Secure: false,
		})
		if err != nil {
			panic(err)
		}
	})
}

func (m *minioImpl) Upload(ctx context.Context, filePath string) (key string, err error) {
	m.init()
	// Upload the test file with FPutObject
	info, err := m.cli.FPutObject(ctx, bucketName, path.Base(filePath), filePath, minio.PutObjectOptions{})
	if err != nil {
		return
	}
	key = info.Key
	return
}

func (m *minioImpl) GetURL(ctx context.Context, key string) (url string, err error) {
	m.init()
	url1, err := m.cli.PresignedGetObject(ctx, bucketName, key, time.Hour*24*3, nil)
	if err != nil {
		return
	}
	url = url1.String()
	return
}
