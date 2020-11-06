package apk

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type Repo struct {
	Location string
	Arch     string

	idx *IndexArchive
}

func (r *Repo) Index() (Index, error) {
	if r.idx == nil {
		reader, err := fetch(fetchJoin(r.Location, r.Arch, "APKINDEX.tar.gz"))
		if err != nil {
			return nil, errors.Wrap(err, "fetch APKINDEX.tar.gz")
		}
		defer reader.Close()
		r.idx, err = ParseApkIndexArchive(reader)
		if err != nil {
			return nil, errors.Wrap(err, "read APKINDEX.tar.gz")
		}
	}
	return r.idx.Index, nil
}
func (r *Repo) Download(e IndexEntry) (io.ReadCloser, error) {
	return fetch(fetchJoin(r.Location, r.Arch, fmt.Sprintf("%s-%s.apk", e.Name, e.Version)))
}
