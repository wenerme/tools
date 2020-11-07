package apk

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

type ArchiveEntryProcessor func(ctx context.Context, h *tar.Header, r io.Reader) error

type ProcessArchiveOptions struct {
	Processors []ArchiveEntryProcessor
	Ungzipped  bool
	Context    context.Context
}

func ProcessArchive(r io.Reader, o *ProcessArchiveOptions) error {
	if o == nil {
		o = &ProcessArchiveOptions{}
	}
	if o.Context == nil {
		o.Context = context.Background()
	}
	var err error
	if !o.Ungzipped {
		r, err = gzip.NewReader(r)
		if err != nil {
			return err
		}
	}
	t := tar.NewReader(r)

	for {
		h, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ctl, r := newBufReader(t)

		for _, f := range o.Processors {
			err = f(o.Context, h, r)
			if err != nil {
				return errors.Wrapf(err, "process %q", h.Name)
			}
			ctl.Reset()
		}
	}
	return nil
}

type bufReader struct {
	r     io.Reader
	all   bool
	rd    int
	buf   *bytes.Buffer
	bytes []byte
	rr    *io.Reader
}

func newBufReader(r io.Reader) (*bufReader, io.Reader) {
	buf := &bytes.Buffer{}
	var rr io.Reader
	lr := &bufReader{
		r:   io.TeeReader(r, buf),
		buf: buf,
		rr:  &rr,
	}
	*lr.rr = lr
	return lr, rr
}
func (r *bufReader) Read(p []byte) (n int, err error) {
	if r.r == nil {
		v := *r.rr
		return v.Read(p)
	}
	n, err = r.r.Read(p)
	if err == io.EOF {
		r.all = true
	}
	r.rd += n
	return n, err
}
func (r *bufReader) Reset() {
	if r.rd == 0 {
		return
	}
	if !r.all {
		_, err := io.Copy(ioutil.Discard, r)
		if err != nil {
			panic(err)
		}
	}
	// switch reader
	if r.bytes == nil {
		r.bytes = r.buf.Bytes()
		r.buf = nil
		r.r = nil
	}
	*r.rr = bytes.NewReader(r.bytes)
}
