package scel

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

const (
	OffsetPinyin  = 0x1540
	OffsetChinese = 0x2628
)

var (
	magic       = []byte{0x40, 0x15, 0x00, 0x00, 0x44, 0x43, 0x53, 0x01, 0x01, 0x00, 0x00, 0x00}
	pinyinMagic = []byte{0x9D, 0x01, 0x00, 0x00}
)

type LoadOptions struct {
	SkipExt           bool
	ResolveWordPinyin bool
}

func LoadFile(fn string, o *LoadOptions) (*Dict, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return LoadBytes(b, o)
}
func LoadBytes(b []byte, o *LoadOptions) (*Dict, error) {
	if o == nil {
		o = &LoadOptions{}
	}
	if len(b) <= OffsetPinyin+len(pinyinMagic) {
		return nil, errors.New("invalid data")
	}
	if !bytes.HasPrefix(b, magic) || !bytes.Equal(b[OffsetPinyin:OffsetPinyin+len(pinyinMagic)], pinyinMagic) {
		return nil, errors.New("invalid data format")
	}
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	decoder := utf16.NewDecoder()
	r := &reader{
		bytes:   b,
		decoder: decoder,
	}

	scel := &Dict{
		Info: Info{
			Name:        r.str(0x130, 0x338),
			Type:        r.str(0x338, 0x540),
			Description: r.str(0x540, 0xD40),
			Example:     r.str(0xD40, OffsetPinyin),
		},
	}

	// pinyin
	pyRemap := make(map[uint16]int)
	pys := make([]string, 0)
	{
		b := b[OffsetPinyin+len(pinyinMagic) : OffsetChinese]

		for len(b) > 0 {
			idx := binary.LittleEndian.Uint16(b)
			l := binary.LittleEndian.Uint16(b[2:])
			b = b[4:]
			s := r.readString(b[:l])
			b = b[l:]

			pyRemap[idx] = len(pyRemap)
			pys = append(pys, s)
		}
	}

	// word
	var words []Word
	{
		b := b[OffsetChinese:]
		for len(b) > 0 {
			w := Word{}
			// 同音词
			same := int(binary.LittleEndian.Uint16(b))
			b = b[2:]

			// 拼音
			pyLen := int(binary.LittleEndian.Uint16(b))
			b = b[2:]
			// 2 per py, pyLen/2
			for i := 0; i < pyLen/2; i++ {
				w.PinyinIndex = append(w.PinyinIndex, pyRemap[binary.LittleEndian.Uint16(b[i*2:])])
			}
			if o.ResolveWordPinyin {
				w.Pinyin = make([]string, len(w.PinyinIndex))
				for i, idx := range w.PinyinIndex {
					w.Pinyin[i] = pys[idx]
				}
			}

			b = b[pyLen:]
			for i := 0; i < same; i++ {
				// 词组
				wordLen := int(binary.LittleEndian.Uint16(b))
				b = b[2:]
				word := r.readString(b[:wordLen])

				b = b[wordLen:]

				// 扩展
				extLen := int(binary.LittleEndian.Uint16(b))
				b = b[2:]
				ext := b[:extLen]

				b = b[extLen:]

				w.Words = append(w.Words, word)
				if !o.SkipExt {
					w.Exts = append(w.Exts, ext)
				}
			}
			words = append(words, w)
		}
	}
	scel.Words = words
	scel.PinyinIndex = pys
	return scel, nil
}

type reader struct {
	bytes   []byte
	decoder *encoding.Decoder
}

func (r *reader) str(a, b int) string {
	return r.readString(r.bytes[a:b])
}

func (r *reader) readString(b []byte) string {
	i := 0
	for ; i < len(b); i += 2 {
		if b[i] == 0 && b[i+1] == 0 {
			break
		}
	}
	if i > 0 {
		dst, err := r.decoder.Bytes(b[:i])
		defer r.decoder.Reset()
		if err != nil {
			panic(err)
		}
		return string(dst)
	}
	return ""
}
