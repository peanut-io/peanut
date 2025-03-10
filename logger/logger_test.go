package logger

import (
	"github.com/peanut-io/peanut/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	err := config.LoadPath("./test/configs")
	assert.NoError(t, err)

	logCfg := &Config{}
	err = config.ScanFrom(logCfg, "logs")
	assert.NoError(t, err)
	assert.NotNil(t, logCfg.File)
	assert.False(t, logCfg.File.Compress)
}
