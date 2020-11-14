package stardict

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func readInfo(reader io.Reader) (info *DictInfo, err error) {
	info = &DictInfo{}
	scanner := bufio.NewScanner(reader)

	if !scanner.Scan() {
		return nil, errors.New("info magic header not found")
	}
	info.Header = scanner.Text()

	for scanner.Scan() {
		err = scanner.Err()
		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			return nil, err
		}
		line := scanner.Text()
		i := strings.IndexRune(line, '=')
		if i < 1 {
			return nil, errors.New("invalid info line: " + line)
		}

		key := strings.TrimSpace(line[:i])
		val := strings.TrimSpace(line[i+1:])
		switch key {
		case "version":
			info.Version = val
		case "wordcount":
			info.WordCount, err = strconv.Atoi(val)
		case "synwordcount":
			info.SynWordCount, err = strconv.Atoi(val)
		case "idxfilesize":
			info.IndexFileSize, err = strconv.Atoi(val)
		case "idxoffsetbits":
			info.IndexOffsetBits, err = strconv.Atoi(val)
		case "bookname":
			info.BookName = val
		case "description":
			info.Description = val
		case "date":
			var t time.Time
			t, err = parseLayouts(val, "2006.1.2", "2006-01-02", "Jan 02 2006")
			if err != nil {
				zap.S().With("value", val, "error", err).Warnf("failed to parse date")
				err = nil
			} else {
				info.Date = &t
			}
		case "sametypesequence":
			info.SameTypeSequence = val
		case "dicttype":
			info.DictType = val
		case "author":
			info.Author = val
		case "email":
			info.Email = val
		case "website":
			info.Website = val
		default:
			return nil, errors.New("unknown info line: " + line)
		}
		if err != nil {
			return nil, errors.New("failed to parse info line: " + line)
		}
	}
	if info.IndexOffsetBits == 0 {
		info.IndexOffsetBits = 32
	}
	return
}

func parseLayouts(s string, layouts ...string) (time.Time, error) {
	var t time.Time
	var err error
	for _, v := range layouts {
		t, err = time.Parse(v, s)
		if err != nil {
			continue
		}
	}
	return t, err
}
