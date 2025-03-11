package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/peanut-io/peanut/config"
)

var storage Storage
var storageOnce sync.Once
var storageConfigOnce sync.Once
var storages map[string]func(*Config) Storage

const (
	Minio = "minio"
)

type Storage interface {
	BucketName() string
	Upload(ctx context.Context, src string, dst string) error
	Download(ctx context.Context, src string, dst string) error
	Read(ctx context.Context, src string) ([]byte, error)
	Write(ctx context.Context, data []byte, dst string) error
	WriteString(ctx context.Context, data string, dst string) error
}

func RegisterStorage(key string, storage func(*Config) Storage) {
	storageOnce.Do(func() {
		storages = make(map[string]func(*Config) Storage)
	})
	if _, ok := storages[key]; !ok {
		storages[key] = storage
	}
}

func GetStorage() Storage {
	storageConfigOnce.Do(func() {
		cfg := &Config{}
		if err := config.ScanFrom(cfg, "storage"); err != nil {
			panic(fmt.Sprintf("failed to get the server config, error: %s", err.Error()))
		} else {
			if caller, ok := storages[cfg.Vendor]; !ok {
				panic(fmt.Sprintf("vendor is not supported, vendor: %s", cfg.Vendor))
			} else {
				storage = caller(cfg)
			}
		}
	})
	return storage
}
