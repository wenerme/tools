package apk

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func fetchString(s string) (string, error) {
	r, err := fetch(s)
	if err != nil {
		return "", err
	}
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func fetch(s string) (io.ReadCloser, error) {
	r, err := http.Get(s) //nolint: gosec,bodyclose
	if err != nil {
		return nil, err
	}
	return r.Body, nil
}

func fetchJoin(args ...string) string {
	return strings.Join(args, "/")
}
