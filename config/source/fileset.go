package source

import (
	"io"
	"os"
)

type FileSet struct {
	Data   map[string]any
	Format string
	Path   string
	Opts   *Options
}

func (fs *FileSet) Read() error {
	f, err := os.Open(fs.Path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	tmpMap := make(map[string]any)
	_ = fs.Opts.Encoder.Decode(data, &tmpMap)
	fs.Data = tmpMap
	return nil
}

func (fs *FileSet) String() string {
	return fs.Path
}

func NewSource(path string, format string) Source {
	options := NewOptions(format)
	return &FileSet{Opts: options, Path: path, Format: format}
}
