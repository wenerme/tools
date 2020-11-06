package apk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	AliyunMirror Mirror = "https://mirrors.aliyun.com/alpine"
	TunaMirror   Mirror = "https://mirrors.tuna.tsinghua.edu.cn/alpine"
)

func TestMirror(t *testing.T) {
	_, err := AliyunMirror.LastUpdated()
	assert.NoError(t, err)

	{
		mirrors, err := TunaMirror.Mirrors()
		assert.NoError(t, err)
		assert.NotEmpty(t, mirrors)
	}

	r := AliyunMirror.Repo("edge", "main", "x86_64")
	idx, err := r.Index()
	assert.NoError(t, err)
	assert.NotNil(t, idx[0])
	assert.NotEmpty(t, idx[0].Name)
}
