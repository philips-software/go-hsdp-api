package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHsdpToml(t *testing.T) {
	info, err := hsdpToml()
	if !assert.Nil(t, err) {
		return
	}
	fileInfo := info.info
	assert.Equal(t, "hsdp.toml", fileInfo.Name())
	assert.Less(t, time.Now().Unix(), fileInfo.ModTime().AddDate(0, 0, 180).Unix())
	assert.False(t, fileInfo.IsDir())
	assert.Less(t, int64(0), fileInfo.Size())
	assert.Equal(t, os.FileMode(420), fileInfo.Mode())
	assert.Nil(t, fileInfo.Sys())
	if !assert.Nil(t, err) {
		return
	}
	name := info.info.Name()
	if !assert.Nil(t, err) {
		return
	}
	_ = MustAsset(name)
	_, err = AssetInfo(name)
	if !assert.Nil(t, err) {
		return
	}
	_, err = AssetDir("foo")
	if !assert.NotNil(t, err) {
		return
	}
	names := AssetNames()
	assert.Equal(t, 1, len(names))
}
