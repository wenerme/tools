package apk

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type IndexArchive struct {
	Sign        []byte
	SignName    string
	Description string
	Index       Index
}

func ParseApkIndexArchive(r io.Reader) (*IndexArchive, error) {
	g, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	t := tar.NewReader(g)
	a := &IndexArchive{}
	for {
		h, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if h.Typeflag != tar.TypeReg {
			continue
		}

		switch h.Name {
		case "APKINDEX":
			b, err := ioutil.ReadAll(t)
			if err != nil {
				return nil, errors.Wrap(err, "read APKINDEX")
			}
			idx := ParseApkIndexList(string(b))
			var i []IndexEntry
			for _, v := range idx {
				e, err := ParseIndexEntry(v)
				if err != nil {
					return nil, errors.Wrapf(err, "parse APKINDEX %q", v["P"])
				}
				i = append(i, e)
			}
			a.Index = i
		case "DESCRIPTION":
			b, err := ioutil.ReadAll(t)
			if err != nil {
				return nil, errors.Wrap(err, "read DESCRIPTION")
			}
			a.Description = string(b)
		}
		if strings.HasPrefix(h.Name, ".SIGN.") {
			a.SignName = h.Name
			b, err := ioutil.ReadAll(t)
			if err != nil {
				return nil, errors.Wrapf(err, "read %q", h.Name)
			}
			a.Sign = b
		}
	}
	return a, nil
}
