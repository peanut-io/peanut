package logger

type Config struct {
	Level        string                 `json:"level" yaml:"level"`
	Encode       string                 `json:"encode" yaml:"encode"`
	LevelPort    int                    `json:"levelPort" yaml:"levelPort"`
	LevelPattern string                 `json:"levelPattern" yaml:"levelPattern"`
	Output       string                 `json:"output" yaml:"output"`
	InitFields   map[string]interface{} `json:"initFields" yaml:"initFields"`
	File         FileSinkConfig         `json:"file" yaml:"file"`
}

type FileSinkConfig struct {
	Path       string `json:"path" yaml:"path"`
	MaxSize    int    `json:"maxSize" yaml:"maxSize"` // megabytes
	MaxBackups int    `json:"maxBackups" yaml:"maxBackups"`
	MaxAge     int    `json:"maxAge" yaml:"maxAge"` // days
	Encode     string `json:"encode"  yaml:"encode"`
	Compress   bool   `json:"compress" yaml:"compress"`
}

var defaultCfg = &Config{
	Level:  "debug",
	Encode: "console",
	Output: "console",
	File: FileSinkConfig{
		Path:       "./logs/app.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Encode:     "json",
		Compress:   false,
	},
}
