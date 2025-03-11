package storage

import (
	"context"
	"github.com/peanut-io/peanut/storage/core"
	_ "github.com/peanut-io/peanut/storage/minio"
)

func Upload(ctx context.Context, src, dst string) error {
	return core.GetStorage().Upload(ctx, src, dst)
}

func Download(ctx context.Context, src, dst string) error {
	return core.GetStorage().Download(ctx, src, dst)
}

func Read(ctx context.Context, src string) ([]byte, error) {
	return core.GetStorage().Read(ctx, src)
}

func Write(ctx context.Context, data []byte, dst string) error {
	return core.GetStorage().Write(ctx, data, dst)
}

func WriteString(ctx context.Context, data string, dst string) error {
	return core.GetStorage().WriteString(ctx, data, dst)
}
