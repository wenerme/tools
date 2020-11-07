package apk

import (
	"bytes"
	"crypto/md5"  // nolint: gosec
	"crypto/sha1" // nolint: gosec
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
)

type ChecksumType int

const (
	ChecksumNone = 0
	ChecksumMd5  = 16
	ChecksumSha1 = 20
)

func (s ChecksumType) String() string {
	switch s {
	case ChecksumNone:
		return "none"
	case ChecksumSha1:
		return "sha1"
	case ChecksumMd5:
		return "md5"
	}
	return fmt.Sprintf("unknown(%v)", int(s))
}

type Checksum struct {
	Type ChecksumType
	Sum  []byte
}

func (s Checksum) CheckBytes(b []byte) error {
	var r []byte
	switch s.Type {
	case ChecksumNone:
		return nil
	case ChecksumSha1:
		v := sha1.Sum(b) // nolint: gosec
		r = v[:]
	case ChecksumMd5:
		v := md5.Sum(b) // nolint: gosec
		r = v[:]
	default:
		return fmt.Errorf("unexpected checksum type %q", s.Type)
	}
	if !bytes.Equal(s.Sum, r) {
		return fmt.Errorf("%s checksum not match", s.Type)
	}
	return nil
}
func (s Checksum) Check(r io.Reader) error {
	var h hash.Hash
	switch s.Type {
	case ChecksumNone:
		return nil
	case ChecksumSha1:
		h = sha1.New() // nolint: gosec
	case ChecksumMd5:
		h = md5.New() // nolint: gosec
	default:
		return fmt.Errorf("unexpected checksum type %q", s.Type)
	}

	if _, err := io.Copy(h, r); err != nil {
		return err
	}
	b := h.Sum(nil)
	if !bytes.Equal(s.Sum, b) {
		return fmt.Errorf("%s checksum not match", s.Type)
	}
	return nil
}

func ParseChecksum(s string) (Checksum, error) {
	c := Checksum{
		Type: ChecksumNone,
	}
	if s == "" {
		return c, nil
	}
	if len(s) < 2 {
		return c, fmt.Errorf("ParseChecksum: invalid checksum size %v", len(s))
	}
	enc, typ := s[0], s[1]
	var err error
	switch enc {
	case 'X':
		c.Sum, err = hex.DecodeString(s[2:])
	case 'Q':
		c.Sum, err = base64.StdEncoding.DecodeString(s[2:])
	}
	if err != nil {
		return c, err
	}
	switch typ {
	case '1':
		c.Type = ChecksumSha1
	default:
		c.Type = ChecksumMd5
	}
	return c, nil
}
