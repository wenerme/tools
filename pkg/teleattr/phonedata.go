package teleattr

import (
	"errors"
	"sort"
	"strconv"
)

type PhoneData struct {
	Version string
	Index   []Index
	Records []Record
}

var ErrInvalidNumber = errors.New("invalid phone number")
var ErrNotFound = errors.New("phone number attribution not found")

func (p *PhoneData) Search(num string) (*Index, error) {
	if len(num) < 7 {
		return nil, ErrInvalidNumber
	}
	pn, err := strconv.Atoi(num[:7])
	if err != nil {
		return nil, ErrInvalidNumber
	}
	i := sort.Search(len(p.Index), func(i int) bool {
		return p.Index[i].Prefix >= pn
	})

	if i < len(p.Index) {
		return &p.Index[i], nil
	}

	return nil, ErrNotFound
}
