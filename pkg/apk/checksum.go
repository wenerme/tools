package apk

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
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
	Data []byte
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
		c.Data, err = hex.DecodeString(s[2:])
	case 'Q':
		c.Data, err = base64.StdEncoding.DecodeString(s[2:])
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
