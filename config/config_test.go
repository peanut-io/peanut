package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPath(t *testing.T) {
	err := LoadPath("./test")
	assert.NoError(t, err)

	var v1 int
	err = ScanFrom(&v1, "TestConfig.value1")
	assert.NoError(t, err)
	assert.Equal(t, 1, v1)

	var v2 string
	err = ScanFrom(&v2, "TestConfig.value2")
	assert.NoError(t, err)
	assert.Equal(t, "a", v2)

	var v3 bool
	err = ScanFrom(&v3, "TestConfig.value3")
	assert.NoError(t, err)
	assert.Equal(t, true, v3)

	var v4 float64
	err = ScanFrom(&v4, "TestConfig.value4")
	assert.NoError(t, err)
	assert.Equal(t, 3.14, v4)
}

func TestLoadFile(t *testing.T) {
	err := LoadFile("./test/test.yaml")
	assert.NoError(t, err)

	v := ""
	err = ScanFrom(&v, "TestConfig.value2")
	assert.NoError(t, err)
	assert.Equal(t, "a", v)
	v2 := false
	err = ScanFrom(&v2, "TestConfig.value3")
	assert.NoError(t, err)
	assert.Equal(t, true, v2)
}

func TestSet(t *testing.T) {
	err := LoadFile("./test/test.yaml")
	assert.NoError(t, err)

	v := ""
	err = ScanFrom(&v, "TestConfig.value2")
	assert.NoError(t, err)
	assert.Equal(t, "a", v)

	Set("TestConfig.value2", "b")
	err = ScanFrom(&v, "TestConfig.value3")
	assert.NoError(t, err)
	assert.Equal(t, "b", v)
}
