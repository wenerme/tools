package scel_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/wenerme/tools/pkg/scel"
)

func TestLoad(t *testing.T) {
	s, err := scel.LoadFile("./testdata/全国省市区县地名.scel", &scel.LoadOptions{})
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, s.Words[0].Exts)
	assert.Nil(t, s.Words[0].Pinyin)

	s, err = scel.LoadFile("./testdata/全国省市区县地名.scel", &scel.LoadOptions{
		ResolveWordPinyin: true,
		SkipExt:           true,
	})
	if !assert.NoError(t, err) {
		return
	}

	assert.Nil(t, s.Words[0].Exts)
	assert.NotNil(t, s.Words[0].Pinyin)

	spew.Dump(s.Info)
	spew.Dump(len(s.Words))
	spew.Dump(s.Words[100])

	word, err := s.Search("shang", "hai")
	assert.NoError(t, err)
	assert.Equal(t, "上海", word.Words[0])
}

func TestInvalidFile(t *testing.T) {
	_, err := scel.LoadFile("./testdata/全国省市区县地名", nil)
	assert.Error(t, err)
	_, err = scel.LoadFile("./load.go", nil)
	assert.Error(t, err)
}
