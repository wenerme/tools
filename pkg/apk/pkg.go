package apk

import (
	"archive/tar"
	"context"
	"encoding/hex"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type PackageInfo struct {
	Arch             string
	BuildDate        time.Time
	Commit           string
	DataHash         int
	Depends          []string
	InstallIf        []string
	License          string
	Maintainer       string
	Origin           string
	Packager         int
	Description      string
	Name             string
	Version          string
	ProviderPriority int
	Provides         []string
	Replaces         []string
	ReplacesPriority int
	Size             int
	Triggers         int
	URL              string
}

func ParsePackageInfoMap(s string) map[string][]string {
	m := make(map[string][]string)
	for _, line := range strings.Split(s, "\n") {
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		i := strings.IndexByte(line, '=')
		if i < 0 {
			continue
		}
		k := strings.TrimSpace(line[:i])
		v := strings.TrimSpace(line[i+1:])
		m[k] = append(m[k], v)
	}
	return m
}

func ParsePackageInfo(info map[string][]string) (PackageInfo, error) {
	m := mapSlice(info)
	var er []error

	pkg := PackageInfo{
		Arch:             m.First("arch"),
		BuildDate:        parseUnixTime(m.First("builddate"), "builddate", &er),
		Commit:           m.First("commit"),
		Depends:          m["depend"],
		InstallIf:        m["install_if"],
		License:          m.First("license"),
		Maintainer:       m.First("maintainer"),
		Origin:           m.First("origin"),
		Packager:         0,
		Description:      m.First("pkgdesc"),
		Name:             m.First("pkgname"),
		Version:          m.First("pkgver"),
		ProviderPriority: parseInt(m.First("provider_priority"), "provider_priority", &er),
		Provides:         m["provide"],
		Replaces:         m["replace"],
		ReplacesPriority: parseInt(m.First("replaces_priority"), "replaces_priority", &er),
		Size:             parseInt(m.First("size"), "size", &er),
		Triggers:         0,
		URL:              m.First("url"),
	}

	return pkg, errorSlice(er)
}

type ParsePackageArchiveOptions struct {
	PackageInfo *PackageInfo
	Checksum    bool
	Processors  []ArchiveEntryProcessor
	Context     context.Context
	Ungzipped   bool
}

// io_archive.c
func ParsePackageArchive(r io.Reader, o *ParsePackageArchiveOptions) error {
	if o == nil {
		o = &ParsePackageArchiveOptions{}
	}
	var p []ArchiveEntryProcessor
	if o.PackageInfo != nil {
		p = append(p, ExtractPackageInfo(o.PackageInfo))
	}
	if o.Checksum {
		p = append(p, ArchiveEntryChecksum)
	}
	if len(o.Processors) > 0 {
		p = append(p, o.Processors...)
	}
	return ProcessArchive(r, &ProcessArchiveOptions{
		Ungzipped:  o.Ungzipped,
		Processors: p,
		Context:    o.Context,
	})
}

func ExtractPackageInfo(info *PackageInfo) ArchiveEntryProcessor {
	return func(ctx context.Context, h *tar.Header, r io.Reader) error {
		if h.Typeflag == tar.TypeReg && h.Name == ".PKGINFO" && info != nil {
			s, err := readAllString(r)
			if err != nil {
				return errors.Wrap(err, "read .PKGINFO")
			}
			*info, err = ParsePackageInfo(ParsePackageInfoMap(s))
			if err != nil {
				return errors.Wrap(err, "parse .PKGINFO")
			}
		}
		return nil
	}
}
func ArchiveEntryChecksum(_ context.Context, h *tar.Header, r io.Reader) error {
	ck := Checksum{
		Type: ChecksumNone,
	}
	var err error

	if v, ok := h.PAXRecords["APK-TOOLS.checksum.SHA1"]; ok {
		ck.Type = ChecksumSha1
		ck.Sum, err = hex.DecodeString(v)
		if err != nil {
			return errors.Wrapf(err, "decode sha1 checksum %q", h.Name)
		}
	}
	if v, ok := h.PAXRecords["APK-TOOLS.checksum.MD5"]; ok {
		ck.Type = ChecksumSha1
		ck.Sum, err = hex.DecodeString(v)
		if err != nil {
			return errors.Wrapf(err, "decode sha1 checksum %q", h.Name)
		}
	}
	if ck.Type != ChecksumNone {
		if err := ck.Check(r); err != nil {
			return errors.Wrapf(err, "%q checksum failed", h.Name)
		}
	}
	return nil
}

type mapSlice map[string][]string

func (m mapSlice) First(n string) string {
	v := m[n]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}
