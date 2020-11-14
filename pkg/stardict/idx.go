package stardict

import (
	"bufio"
	"encoding/binary"
	"io"
)

func readIndexEntries(bits int, reader io.Reader) (entries []DictIndexEntry, err error) {
	r := bufio.NewReader(reader)
	buf := make([]byte, 4)
	var b []byte
	for {
		entry := DictIndexEntry{}

		// word
		b, err = r.ReadBytes(0)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}
		entry.Word = string(b[:len(b)-1]) // trim \x00

		// offset
		if bits == 32 {
			_, err = io.ReadFull(r, buf[:4])
			entry.Offset = int64(binary.BigEndian.Uint32(buf))
		} else {
			_, err = io.ReadFull(r, buf)
			entry.Offset = int64(binary.BigEndian.Uint64(buf))
		}
		if err != nil {
			return nil, err
		}

		// size
		_, err = io.ReadFull(r, buf[:4])
		if err != nil {
			return nil, err
		}
		entry.Size = int(binary.BigEndian.Uint32(buf))
		entries = append(entries, entry)
	}
	return
}
