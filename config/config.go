package config

import (
	"strings"
)

var defaultConfig = newConfig()

func newConfig() *config {
	cfg := &config{}
	cfg.Load(defaultSources()...)
	return cfg
}

func Scan(key interface{}) error {
	return defaultConfig.Scan(key)
}

func ScanFrom(v interface{}, key string) error {
	val := Get(key)
	return val.Scan(v)
}

func Get(path ...string) Values {
	return defaultConfig.Get(normalizePath(path...)...)
}

func Set(key string, value interface{}) {
	defaultConfig.Set(normalizePath(key), value)
}

func LoadPath(path string) error {
	return defaultConfig.Load(newSourcesFromDir(path)...)
}

func LoadFile(filePath string) error {
	return defaultConfig.Load(newSourcesFromFile(filePath))
}

func Reload() error {
	return defaultConfig.Load(defaultSources()...)
}

func normalizePath(path ...string) []string {
	var segments []string
	for _, p := range path {
		s := strings.Split(p, ".")
		segments = append(segments, s...)
	}
	return segments
}
