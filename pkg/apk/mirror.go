package apk

import (
	"strconv"
	"strings"
	"time"
)

type Mirror string

func (m Mirror) Mirrors() ([]string, error) {
	s, err := fetchString(fetchJoin(string(m), "MIRRORS.txt"))
	if err != nil {
		return nil, err
	}
	s = strings.TrimSpace(s)
	return strings.Split(s, "\n"), nil
}

func (m Mirror) LastUpdated() (time.Time, error) {
	s, err := fetchString(fetchJoin(string(m), "last-updated"))
	if err != nil {
		return time.Time{}, err
	}
	s = strings.TrimSpace(s)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0), nil
}

func (m Mirror) Repo(ver string, repo string, arch string) Repo {
	return Repo{
		Location: fetchJoin(string(m), ver, repo),
		Arch:     arch,
	}
}
