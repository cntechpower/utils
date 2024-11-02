package oss

import (
	"context"
)

var Impl Storage

type Storage interface {
	Upload(ctx context.Context, filePath string) (key string, err error)
	GetURL(ctx context.Context, key string) (url string, err error)
}
