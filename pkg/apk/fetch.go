package apk

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func fetchString(s string) (string, error) {
	r, err := fetch(s)
	if err != nil {
		return "", err
	}
	defer r.Close()
	return readAllString(r)
}

func fetch(s string) (io.ReadCloser, error) {
	if s[0] == '/' {
		return os.Open(s)
	}
	r, err := http.Get(s) //nolint: gosec,bodyclose
	if err != nil {
		return nil, err
	}
	return r.Body, nil
}

func fetchJoin(args ...string) string {
	return strings.Join(args, "/")
}

func readAllString(r io.Reader) (string, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
