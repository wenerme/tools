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

func (r *Repo) IndexArchive() (IndexArchive, error) {
	if r.idx == nil {
		reader, err := fetch(fetchJoin(r.Location, r.Arch, "APKINDEX.tar.gz"))
		if err != nil {
			return IndexArchive{}, errors.Wrap(err, "fetch APKINDEX.tar.gz")
		}
		defer reader.Close()
		r.idx, err = ParseApkIndexArchive(reader)
		if err != nil {
			return IndexArchive{}, errors.Wrap(err, "read APKINDEX.tar.gz")
		}
	}
	return *r.idx, nil
}
func (r Repo) Index() (Index, error) {
	ar, err := r.IndexArchive()
	return ar.Index, err
}
func (r *Repo) Download(e IndexEntry) (io.ReadCloser, error) {
	return fetch(fetchJoin(r.Location, r.Arch, fmt.Sprintf("%s-%s.apk", e.Name, e.Version)))
}
func (r Repo) String() string {
	return fetchJoin(r.Location, r.Arch)
}
