package stardict

import (
	"archive/tar"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func Open(f string) (*Dict, error) {
	if fi, err := os.Stat(f); err != nil {
		return nil, err
	} else if fi.IsDir() {
		return openDir(f)
	}
	return openArchive(f)
}

func openArchive(fn string) (*Dict, error) {
	ext := filepath.Ext(fn)
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var r io.Reader
	{
		var err error
		r = f
		switch ext {
		case ".bz2":
			r = bzip2.NewReader(r)
		case ".gz":
			r, err = gzip.NewReader(r)
		case ".tar":
		default:
			err = errors.New("unrecognized file format: " + ext)
		}
		if err != nil {
			return nil, err
		}
	}
	dict := &Dict{}
	{
		var idxData []byte
		var dictData []byte
		err := ProcessArchive(r, &ProcessArchiveOptions{
			Processors: []ArchiveEntryProcessor{
				func(ctx context.Context, h *tar.Header, r io.Reader) (err error) {
					if strings.HasSuffix(h.Name, ".ifo") {
						dict.Info, err = readInfo(r)
					}
					return
				},
				func(ctx context.Context, h *tar.Header, r io.Reader) (err error) {
					if strings.HasSuffix(h.Name, ".idx") {
						idxData, err = ioutil.ReadAll(r)
					}
					return
				},
				func(ctx context.Context, h *tar.Header, r io.Reader) (err error) {
					if strings.HasSuffix(h.Name, ".idx.gz") {
						r, err = gzip.NewReader(r)
						if err == nil {
							idxData, err = ioutil.ReadAll(r)
						}
					}
					return
				},

				func(ctx context.Context, h *tar.Header, r io.Reader) (err error) {
					if strings.HasSuffix(h.Name, ".dict") {
						dictData, err = ioutil.ReadAll(r)
					}
					return
				},
				func(ctx context.Context, h *tar.Header, r io.Reader) (err error) {
					if strings.HasSuffix(h.Name, ".dict.dz") {
						r, err = gzip.NewReader(r)
						if err == nil {
							dictData, err = ioutil.ReadAll(r)
						}
					}
					return
				},
			},
		})
		if err != nil {
			return dict, err
		}
		if idxData == nil {
			return nil, errors.New("index not found")
		}
		if dictData == nil {
			return nil, errors.New("dict not found")
		}
		dict.Index, err = readIndexEntries(dict.Info.IndexOffsetBits, bytes.NewReader(idxData))
		if err != nil {
			return dict, err
		}
		dict.r = bytes.NewReader(dictData)
	}
	return dict, nil
}
func openDir(dir string) (*Dict, error) {
	m, err := filepath.Glob(fmt.Sprintf("%s/*.ifo", dir))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to detect dict.ifo")
	}
	if len(m) == 0 {
		return nil, errors.New("find non .ifo")
	}
	if len(m) > 1 {
		return nil, errors.New("find more than one .ifo")
	}
	fn := path.Base(m[0])
	name := fn[0:(len(fn) - 4)]

	dict := &Dict{}
	{
		f, err := os.Open(m[0])
		if err != nil {
			return nil, err
		}
		defer f.Close()
		dictInfo, err := readInfo(f)
		if err != nil {
			return nil, err
		}
		dict.Info = dictInfo
	}
	{
		var r io.ReadCloser
		var gz bool
		f, err := os.Open(fmt.Sprintf("%s/%s.idx", dir, name))
		if os.IsNotExist(err) {
			f, err = os.Open(fmt.Sprintf("%s/%s.idx.gz", dir, name))
			gz = true
		}
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
		if gz {
			r, err = gzip.NewReader(f)
			if err != nil {
				return nil, err
			}
		}

		idx, err := readIndexEntries(dict.Info.IndexOffsetBits, r)
		if err != nil {
			return nil, err
		}
		dict.Index = idx
	}

	{
		var r io.ReadCloser
		var gz bool
		f, err := os.Open(fmt.Sprintf("%s/%s.dict", dir, name))
		if os.IsNotExist(err) {
			f, err = os.Open(fmt.Sprintf("%s/%s.dict.gz", dir, name))
			gz = true
		}
		if os.IsNotExist(err) {
			f, err = os.Open(fmt.Sprintf("%s/%s.dict.dz", dir, name))
			gz = true
		}
		if err != nil {
			return nil, err
		}
		r = f
		if gz {
			r, err = gzip.NewReader(f)
			if err != nil {
				return nil, err
			}
		}

		if rs, ok := r.(io.ReadSeeker); ok {
			dict.r = rs
			dict.closers = append(dict.closers, r)
		} else {
			defer r.Close()
			b, err := ioutil.ReadAll(r)
			if err != nil {
				return nil, err
			}
			buf := bytes.NewReader(b)
			dict.r = buf
		}
	}

	return dict, nil
}
