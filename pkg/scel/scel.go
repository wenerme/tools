package scel

import (
	"fmt"
	"sort"
)

type Dict struct {
	Info        Info
	PinyinIndex []string
	Words       []Word // Words is sorted by Word.PinyinIndex
}
type Info struct {
	Name        string
	Type        string
	Description string
	Example     string
}
type Word struct {
	PinyinIndex []int
	Pinyin      []string
	Words       []string
	Exts        [][]byte
}

func (d *Dict) Search(pys ...string) (*Word, error) {
	idx := make([]int, len(pys))
	pinyinIndex := d.PinyinIndex
	for i, v := range pys {
		id := sort.SearchStrings(pinyinIndex, v)
		if id < len(pinyinIndex) && pinyinIndex[id] == v {
			idx[i] = id
		} else {
			return nil, fmt.Errorf("cannot resolve pinyin %q", v)
		}
	}

	words := d.Words
	found := sort.Search(len(words), func(i int) bool {
		return compare(words[i].PinyinIndex, idx) >= 0
	})
	if found < len(words) && compare(words[found].PinyinIndex, idx) == 0 {
		return &words[found], nil
	}
	return nil, nil
}

func compare(a, b []int) int {
	if len(a) >= len(b) {
		for i, vb := range b {
			va := a[i]
			if va == vb {
				continue
			}
			return va - vb
		}
		return 0
	}
	for i, va := range a {
		vb := b[i]
		if va == vb {
			continue
		}
		return va - vb
	}
	return 0
}
