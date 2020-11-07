package apk_test

import (
	"archive/tar"
	"context"
	"crypto/sha1" // nolint: gosec
	"encoding/hex"
	"fmt"
	"io"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/wenerme/tools/pkg/apk"

	"github.com/stretchr/testify/assert"
)

const (
	AliyunMirror apk.Mirror = "https://mirrors.aliyun.com/alpine"
	TunaMirror   apk.Mirror = "https://mirrors.tuna.tsinghua.edu.cn/alpine"
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

	reader, err := r.Download(idx[0])
	assert.NoError(t, err)
	defer reader.Close()

	opts := &apk.ParsePackageArchiveOptions{
		PackageInfo: &apk.PackageInfo{},
		Checksum:    true,
		Processors: []apk.ArchiveEntryProcessor{
			apk.ArchiveEntryChecksum,
			func(ctx context.Context, h *tar.Header, r io.Reader) error {
				hash := sha1.New() // nolint: gosec
				if _, err := io.Copy(hash, r); err != nil {
					return err
				}
				fmt.Println(h.Name, h.Typeflag, h.PAXRecords, h.Mode, h.ModTime, "SHA1", hex.EncodeToString(hash.Sum(nil)))
				return nil
			},
		},
	}
	assert.NoError(t, apk.ParsePackageArchive(reader, opts))
	spew.Dump(opts.PackageInfo)
}
