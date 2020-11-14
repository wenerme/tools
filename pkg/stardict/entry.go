package stardict

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
)

type EntryContentReader func(dict *Dict, idx *DictIndexEntry, typ EntryType, buf *bytes.Buffer) (DictEntry, error)

var EntryContentReaders = make(map[EntryType]EntryContentReader)

type DictEntry struct {
	Type EntryType
	Text string
}

func readTextEntry(dict *Dict, idx *DictIndexEntry, typ EntryType, buf *bytes.Buffer) (DictEntry, error) {
	entry := DictEntry{}
	entry.Type = typ
	b, err := buf.ReadBytes(0)
	if err != nil && err != io.EOF {
		return entry, err
	}
	if len(b) > 0 && b[len(b)-1] == 0 {
		entry.Text = string(b[:len(b)-1])
	} else {
		entry.Text = string(b)
	}
	return entry, nil
}

func init() {
	EntryContentReaders[TypeNullTerminalText] = readTextEntry
	EntryContentReaders[TypeEnglishPhonetic] = readTextEntry
	EntryContentReaders[TypeYinBiao] = readTextEntry
	EntryContentReaders[TypeKingsoftXML] = readTextEntry
	EntryContentReaders[TypeHTML] = readTextEntry
	EntryContentReaders[TypeXdxfMarkup] = readTextEntry
	EntryContentReaders[TypePangoText] = readTextEntry
}

func readDictContent(dict *Dict, idx *DictIndexEntry) error {
	// loaded
	if len(idx.Contents) > 0 {
		return nil
	}

	r := dict.r
	if _, err := r.Seek(idx.Offset, io.SeekStart); err != nil {
		return err
	}
	var contents []DictEntry

	b := make([]byte, idx.Size)
	n, err := io.ReadFull(r, b)
	if err != nil && err != io.EOF {
		return errors.Wrap(err, "failed to read entry data")
	} else if n != idx.Size {
		return errors.Errorf("entry data size not match: %v != %v", n, idx.Size)
	}
	buf := bytes.NewBuffer(b)

	for _, v := range dict.Info.SameTypeSequence {
		typ := EntryType(v)
		reader := EntryContentReaders[typ]
		if reader == nil {
			return errors.New("no reader found for type '" + string(typ) + "'")
		}

		content, err := reader(dict, idx, typ, buf)
		if err != nil {
			return errors.Wrap(err, "failed to read entry")
		}
		contents = append(contents, content)
	}
	idx.Contents = contents
	return nil
}
