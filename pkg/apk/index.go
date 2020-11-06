package apk

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func ParseApkIndexList(s string) []map[string]string {
	idx := make([]map[string]string, 0)
	for _, p := range strings.Split(s, "\n\n") {
		m := make(map[string]string)

		for _, v := range strings.Split(p, "\n") {
			i := strings.IndexByte(v, ':')
			if i < 0 {
				m[v] = ""
			} else {
				m[v[:i]] = v[i+1:]
			}
		}

		idx = append(idx, m)
	}
	return idx
}

type Index []IndexEntry

// https://wiki.alpinelinux.org/wiki/Apk_spec
type IndexEntry struct {
	Arch             string
	Checksum         Checksum
	FilePath         string
	Name             string
	Version          string
	Size             int
	InstallSize      int
	Description      string
	URL              string
	License          string
	Maintainer       string
	Origin           string
	BuildTime        time.Time
	Commit           string
	ProviderPriority int

	Depends   []string
	Provides  []string
	InstallIf []string
}

func ParseIndexEntry(m map[string]string) (IndexEntry, error) {
	var er []error

	// order based on apk-tools apk_pkg_write_index_entry
	e := IndexEntry{
		Checksum: func(s string, msg string, e *[]error) Checksum {
			c, err := ParseChecksum(s)
			if err != nil {
				*e = append(*e, errors.Wrap(err, msg))
			}
			return c
		}(m["C"], "checksum", &er),
		Name:             m["P"],
		Version:          m["V"],
		Arch:             m["A"],
		Size:             parseInt(m["S"], "size", &er),
		InstallSize:      parseInt(m["I"], "install size", &er),
		Description:      m["T"],
		URL:              m["U"],
		License:          m["L"],
		Origin:           m["o"],
		Maintainer:       m["m"],
		BuildTime:        parseUnixTime(m["t"], "build time", &er),
		Commit:           m["c"],
		ProviderPriority: parseInt(m["k"], "provider priority", &er),

		Depends:   parseStringSlice(m["D"]),
		Provides:  parseStringSlice(m["p"]),
		InstallIf: parseStringSlice(m["i"]),
	}

	return e, errorSlice(er)
}
func parseStringSlice(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, " ")
}
func parseUnixTime(s string, msg string, e *[]error) time.Time {
	if s == "" {
		return time.Time{}
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		*e = append(*e, errors.Wrap(err, msg))
	}
	return time.Unix(i, 0)
}
func parseInt(s string, msg string, e *[]error) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		*e = append(*e, errors.Wrap(err, msg))
	}
	return i
}
