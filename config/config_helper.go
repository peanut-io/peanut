package config

import (
	"dario.cat/mergo"
	"os"
	"path"
	"path/filepath"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/peanut-io/peanut/config/source"
)

type config struct {
	values Values
}

func defaultSources() []source.Source {
	workDir, _ := os.Getwd()
	dirs := []string{
		filepath.Join(workDir, "conf"),
		filepath.Join(workDir, "configs"),
	}

	if configPath := getEnv("CONFIG_PATH"); len(configPath) > 0 {
		dirs = append([]string{configPath}, dirs...)
	}

	// debug mode in the cmd folder
	if strings.Contains(workDir, "/cmd/") || strings.HasSuffix(workDir, "/cmd") {
		dirs = append(dirs, []string{"../configs", "../../configs"}...)
	}

	var sources []source.Source
	for _, dir := range dirs {
		if sources = newSourcesFromDir(dir); len(sources) > 0 {
			break
		}
	}
	return sources
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return os.Getenv(key)
	}
	return ""
}

func newSourcesFromDir(dir string) []source.Source {
	var sources []source.Source
	files, _ := os.ReadDir(dir)
	for _, f := range files {
		if f.IsDir() {
			ss := newSourcesFromDir(f.Name())
			sources = append(sources, ss...)
		} else {
			segments := strings.Split(f.Name(), ".")
			suffix := ""
			if len(segments) >= 2 {
				suffix = segments[len(segments)-1]
			}
			if !source.SupportedFileSuffixes[suffix] {
				continue
			}
			p := path.Join(dir, f.Name())

			sources = append(sources, source.NewSource(p, suffix))
		}
	}

	return sources
}

func newSourcesFromFile(filePath string) source.Source {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil
	}
	segments := strings.Split(path.Base(filePath), ".")
	suffix := ""
	if len(segments) >= 2 {
		suffix = segments[len(segments)-1]
	}
	if !source.SupportedFileSuffixes[suffix] {
		return nil
	}
	return source.NewSource(filePath, suffix)
}

func (c *config) Load(sources ...source.Source) error {
	data := make(map[string]any)
	for _, src := range sources {
		fileSet := src.(*source.FileSet)
		err := fileSet.Read()
		if err != nil {
			return err
		}
		err = mergo.Merge(&data, fileSet.Data)
		if err != nil {
			return err
		}
	}

	c.values = newValues(data)
	return nil
}

func (c *config) Scan(v interface{}) error {
	b, err := c.values.sj.MarshalJSON()
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(b, v)
}

func (c *config) Get(path ...string) Values {
	return c.values.Get(path...)
}

func (c *config) Set(key []string, value interface{}) {
	c.values.Set(key, value)
}
